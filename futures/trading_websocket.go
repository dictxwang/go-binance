package futures

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const (
	redialTick = 2 * time.Second
	writeWait  = 3 * time.Second
	pongWait   = 25 * time.Second
	PingPeriod = 15 * time.Second
)

type (
	ErrorDetail struct {
		Code int64  `json:"code"`
		Msg  string `json:"msg,omitempty"`
	}

	Basic struct {
		ID     string                 `json:"id"`
		Status int                    `json:"status"`
		Result map[string]interface{} `json:"result,omitempty"`
		Error  *ErrorDetail           `json:"error,omitempty"`
	}

	Error struct {
		ID     string       `json:"id,omitempty"`
		Status int64        `json:"status,omitempty"`
		Error  *ErrorDetail `json:"error,omitempty"`
	}

	LoginResult struct {
		APIKey           string `json:"apiKey"`
		AuthorizedSince  int64  `json:"authorizedSince"`
		ConnectedSince   int64  `json:"connectedSince"`
		ReturnRateLimits bool   `json:"returnRateLimits"`
		ServerTime       int64  `json:"serverTime"`
	}

	LoginResp struct {
		ID     string      `json:"id,omitempty"`
		Status int64       `json:"status,omitempty"`
		Result LoginResult `json:"result,omitempty"`
	}

	OrderResult struct {
		OrderId                 int    `json:"orderId,omitempty"`
		Symbol                  string `json:"symbol,omitempty"`
		Status                  string `json:"status,omitempty"`
		ClientOrderId           string `json:"clientOrderId,omitempty"`
		Price                   string `json:"price,omitempty"`
		AvgPrice                string `json:"avgPrice,omitempty"`
		OrigQty                 string `json:"origQty,omitempty"`
		ExecutedQty             string `json:"executedQty,omitempty"`
		CumQty                  string `json:"cumQty,omitempty"`
		CumQuote                string `json:"cumQuote,omitempty"`
		TimeInForce             string `json:"timeInForce,omitempty"`
		Type                    string `json:"type,omitempty"`
		ReduceOnly              bool   `json:"reduceOnly,omitempty"`
		ClosePosition           bool   `json:"closePosition,omitempty"`
		Side                    string `json:"side,omitempty"`
		PositionSide            string `json:"positionSide,omitempty"`
		StopPrice               string `json:"stopPrice,omitempty"`
		WorkingType             string `json:"workingType,omitempty"`
		PriceProtect            bool   `json:"priceProtect,omitempty"`
		OrigType                string `json:"origType,omitempty"`
		PriceMatch              string `json:"priceMatch,omitempty"`
		SelfTradePreventionMode string `json:"selfTradePreventionMode,omitempty"`
		GoodTillDate            int    `json:"goodTillDate,omitempty"`
		UpdateTime              int64  `json:"updateTime,omitempty"`
		ActivatePrice           string `json:"activatePrice,omitempty"`
		PriceRate               string `json:"priceRate,omitempty"`
	}

	OrderResp struct {
		ID     string      `json:"id,omitempty"`
		Status int64       `json:"status,omitempty"`
		Result OrderResult `json:"result,omitempty"`
	}
)

func (a *Basic) GetResult(k string) (interface{}, bool) {
	v, ok := a.Result[k]
	return v, ok
}

type ClientWs struct {
	url           string
	apiKey        string
	secretKey     string
	privateKey    crypto.PrivateKey
	conn          *websocket.Conn
	closed        bool
	ctx           context.Context
	Cancel        context.CancelFunc
	DoneChan      chan interface{}
	ErrChan       chan *Error
	LoginChan     chan *LoginResp
	OrderRespChan chan *OrderResp
	sendChan      chan []byte
	AuthRequested *time.Time
	Authorized    bool
	LocalIP       string
	lastTransmit  *time.Time
	resolver      *net.Resolver
}

func NewTradingWsClient(ctx context.Context, apiKey, secretKey, localIP string) *ClientWs {
	ctx, cancel := context.WithCancel(ctx)
	c := &ClientWs{
		url:       getTradingWsEndpoint(),
		apiKey:    apiKey,
		secretKey: secretKey,
		conn:      nil,
		closed:    false,
		ctx:       ctx,
		Cancel:    cancel,
		sendChan:  make(chan []byte, 3),
		DoneChan:  make(chan interface{}, 32),
		LocalIP:   localIP,
	}

	privateKey, err := parsePrivateKey(c.secretKey)
	if err != nil {
		fmt.Errorf("failed to parse private key")
		os.Exit(-1)
	}
	c.privateKey = privateKey
	return c
}

func (c *ClientWs) SetResolver(resolver *net.Resolver) {
	c.resolver = resolver
}

func (c *ClientWs) SetChannels(errCh chan *Error, lCh chan *LoginResp, osCh chan *OrderResp) {
	c.ErrChan = errCh
	c.LoginChan = lCh
	c.OrderRespChan = osCh
}

func (c *ClientWs) Send(method string, args map[string]interface{}, extras ...map[string]string) error {
	if method != "session.logon" {
		err := c.Connect()
		if err == nil {
			err = c.WaitForAuthorization()
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	reqID := GetReqID()
	data := map[string]interface{}{
		"id":     reqID,
		"method": method,
		"params": args,
	}

	for _, extra := range extras {
		for k, v := range extra {
			data[k] = v
		}
	}

	j, err := json.Marshal(data)
	if err != nil {
		return err
	}

	c.sendChan <- j
	return nil
}

func (c *ClientWs) Connect() error {
	if c.CheckConnect() {
		return nil
	}

	err := c.dial()
	if err == nil {
		return nil
	}

	ticker := time.NewTicker(redialTick)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err = c.dial()
			if err == nil {
				return nil
			}
		case <-c.ctx.Done():
			return c.handleCancel("connect")
		}
	}
}

// CheckConnect into the server
func (c *ClientWs) CheckConnect() bool {
	if c.conn != nil && !c.closed {
		return true
	}
	return false
}

// WaitForAuthorization waits for the auth response and try to log in if it was needed
func (c *ClientWs) WaitForAuthorization() error {
	if c.Authorized {
		return nil
	}

	if err := c.Login(); err != nil {
		return err
	}

	ticker := time.NewTicker(time.Millisecond * 300)
	defer ticker.Stop()

	for range ticker.C {
		if c.Authorized {
			return nil
		}
	}

	return nil
}

func (c *ClientWs) Login() error {
	if c.Authorized {
		return nil
	}

	if c.AuthRequested != nil && time.Since(*c.AuthRequested).Seconds() < 30 {
		return nil
	}

	now := time.Now()
	c.AuthRequested = &now

	method := "session.logon"
	args := map[string]interface{}{
		"apiKey":    c.apiKey,
		"timestamp": time.Now().UnixMilli(), // use the current time in milliseconds
	}

	payload := makeQueryString(args)
	signature, err := signPayload(payload, c.privateKey)

	if err != nil {
		fmt.Printf("Failed to sign payload: %v", err)
		return err
	}
	args["signature"] = signature

	reqID := GetReqID()
	data := map[string]interface{}{
		"id":     reqID,
		"method": method,
		"params": args,
	}

	j, err := json.Marshal(data)
	if err != nil {
		return err
	}

	c.sendChan <- j
	return nil
}

func (c *ClientWs) dial() error {
	var dialer websocket.Dialer
	if c.LocalIP != "" {
		dialer = websocket.Dialer{
			NetDial: func(network, addr string) (net.Conn, error) {
				localAddr, err := net.ResolveTCPAddr("tcp", c.LocalIP+":0") // 替换为您的出口IP地址
				if err != nil {
					return nil, err
				}
				var d net.Dialer
				if c.resolver == nil {
					d = net.Dialer{
						LocalAddr: localAddr,
						Resolver:  net.DefaultResolver,
					}
				} else {
					d = net.Dialer{
						LocalAddr: localAddr,
						Resolver:  c.resolver,
					}
				}
				return d.Dial(network, addr)
			},
			HandshakeTimeout:  45 * time.Second,
			EnableCompression: false,
		}
	} else {
		dialer = websocket.Dialer{
			Proxy:             http.ProxyFromEnvironment,
			HandshakeTimeout:  45 * time.Second,
			EnableCompression: false,
		}
	}
	conn, res, err := dialer.Dial(c.url, nil)
	conn.SetReadLimit(655350)
	if err != nil {
		var statusCode int
		if res != nil {
			statusCode = res.StatusCode
		}

		return fmt.Errorf("error %d: %w", statusCode, err)
	}
	defer res.Body.Close()

	go func() {
		defer func() {
			// Cleaning the connection with ws
			c.Cancel()
			c.conn.Close()
			c.closed = true
			fmt.Printf("receiver connection closed\n")
		}()
		err := c.receiver()
		if err != nil {
			if !strings.Contains(err.Error(), "operation cancelled: receiver") {
				c.ErrChan <- &Error{
					Error: &ErrorDetail{
						Code: 111,
						Msg:  err.Error(),
					},
				}
			}
			fmt.Printf("receiver error: %v\n", err)
		}
	}()

	go func() {
		defer func() {
			// Cleaning the connection with ws
			c.Cancel()
			c.conn.Close()
			c.closed = true
			fmt.Printf("sender connection closed\n")
		}()
		err := c.sender()
		if err != nil {
			if !strings.Contains(err.Error(), "operation cancelled: sender") {
				c.ErrChan <- &Error{
					Error: &ErrorDetail{
						Code: 111,
						Msg:  err.Error(),
					},
				}
			}
			fmt.Printf("sender error: %v\n", err)
			c.Authorized = false
		}
	}()

	c.conn = conn
	c.closed = false

	return nil
}

func (c *ClientWs) sender() error {
	ticker := time.NewTicker(time.Millisecond * 300)
	defer ticker.Stop()

	for {
		select {
		case data := <-c.sendChan:
			if string(data) == "ping" {
				deadline := time.Now().Add(10 * time.Second)
				err := c.conn.WriteControl(websocket.PingMessage, []byte{}, deadline)
				if err != nil {
					return fmt.Errorf("failed to send ping to conn, error: %w", err)
				}
			} else {
				err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err != nil {
					return fmt.Errorf("failed to set write deadline for ws connection, error: %w", err)
				}

				err = c.conn.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.ClosePolicyViolation) {
						return fmt.Errorf("connection closed, error: %w", err)
					}

					return fmt.Errorf("Failed to send auth request: %v", err)
				}
				//
				//w, err := c.conn.NextWriter(websocket.TextMessage)
				//if err != nil {
				//	return fmt.Errorf("failed to get next writer for ws connection, error: %w", err)
				//}
				//
				//if _, err = w.Write(data); err != nil {
				//	return fmt.Errorf("failed to write data via ws connection, error: %w", err)
				//}
				//
				//if err := w.Close(); err != nil {
				//	return fmt.Errorf("failed to close ws connection, error: %w", err)
				//}
			}

		case <-ticker.C:
			lastTransmit := c.lastTransmit
			if c.conn != nil && (lastTransmit == nil || (lastTransmit != nil && time.Since(*lastTransmit) > PingPeriod)) {
				go func() {
					c.sendChan <- []byte("ping")
				}()
			}

		case <-c.ctx.Done():
			return c.handleCancel("sender")
		}
	}
}
func (c *ClientWs) receiver() error {
	for {
		select {
		case <-c.ctx.Done():
			return c.handleCancel("receiver")
		default:
			mt, data, err := c.conn.ReadMessage()
			if err != nil {
				return fmt.Errorf("failed to read message from ws connection, error: %v\n", err)
			}

			now := time.Now()
			c.lastTransmit = &now

			if mt == websocket.TextMessage && string(data) != "pong" {
				//fmt.Printf("Raw JSON data: %s\n", data)

				// Attempt to unmarshal into Basic struct
				e := &Basic{}
				if err := json.Unmarshal(data, e); err != nil {
					return fmt.Errorf("Failed to unmarshal message from ws, error: %v\n", err)
				}
				go c.process(data, e)
			}
		}
	}
}

func (c *ClientWs) handleCancel(msg string) error {
	go func() {
		c.DoneChan <- msg
	}()

	return fmt.Errorf("operation cancelled: %s", msg)
}

// TODO: break each case into a separate function
func (c *ClientWs) process(data []byte, e *Basic) bool {

	if e.Error != nil {
		e := Error{}
		_ = json.Unmarshal(data, &e)
		go func() {
			if c.ErrChan != nil {
				c.ErrChan <- &e
			}
		}()
		return true
	}

	if e.Result != nil {
		_, ok := e.GetResult("authorizedSince")
		if ok {
			// logon request
			if time.Since(*c.AuthRequested).Seconds() > 30 {
				c.AuthRequested = nil
				_ = c.Login()
				return false
			}

			c.Authorized = true

			e := LoginResp{}
			_ = json.Unmarshal(data, &e)
			go func() {
				if c.LoginChan != nil {
					c.LoginChan <- &e
				}
			}()

			return true
		}

		orderStatus, ok := e.GetResult("status")
		if !ok {
			return false
		} else {
			if orderStatus == "NEW" || orderStatus == "CANCELED" {
				e := OrderResp{}
				_ = json.Unmarshal(data, &e)
				go func() {
					if c.OrderRespChan != nil {
						c.OrderRespChan <- &e
					}
				}()
			} else {
				return false
			}
		}

	}

	return false
}

type WsPlaceOrder struct {
	NewClientOrderId string  `json:"newClientOrderId"`
	Symbol           string  `json:"symbol"`
	Price            float64 `json:"price"`
	Quantity         float64 `json:"quantity"`
	Side             string  `json:"side"`
	Type             string  `json:"type"`
	TimeInForce      string  `json:"timeInForce"`
	PositionSide     string  `json:"positionSide"`
	Timestamp        int64   `json:"timestamp"`
}

type WsCancelOrder struct {
	Symbol            string `json:"symbol"`
	OrderID           int64  `json:"orderId,omitempty"`
	OrigClientOrderId string `json:"origClientOrderId,omitempty"`
	Timestamp         int64  `json:"timestamp"`
}

func (c *ClientWs) PlaceOrder(order *WsPlaceOrder) error {
	args := s2m(order)
	return c.Send("order.place", args)
}

func (c *ClientWs) CancelOrder(order *WsCancelOrder) error {
	args := s2m(order)
	return c.Send("order.cancel", args)
}

func s2m(i interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	j, _ := json.Marshal(i)
	_ = json.Unmarshal(j, &m)

	return m
}

func getTimestampInMS() int64 {
	return time.Now().UnixNano() / 1e6
}

var gClientOrderID = getTimestampInMS()

func GetReqID() string {
	atomic.AddInt64(&gClientOrderID, 1)
	return strconv.FormatInt(atomic.LoadInt64(&gClientOrderID), 10)
}

func makeQueryString(params map[string]interface{}) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var payloadBuilder strings.Builder
	for _, key := range keys {
		payloadBuilder.WriteString(fmt.Sprintf("%s=%v&", key, params[key]))
	}
	return strings.TrimRight(payloadBuilder.String(), "&")
}

func parsePrivateKey(apiSecret string) (crypto.PrivateKey, error) {
	block, _ := pem.Decode([]byte(apiSecret))
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func signPayload(payload string, privateKey crypto.PrivateKey) (string, error) {
	key, ok := privateKey.(crypto.Signer)
	if !ok {
		return "", fmt.Errorf("key does not implement crypto.Signer")
	}

	signature, err := key.Sign(nil, []byte(payload), crypto.Hash(0))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}
