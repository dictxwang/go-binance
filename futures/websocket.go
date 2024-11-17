package futures

import (
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// WsHandler handle raw websocket message
type WsHandler func(message []byte)

// ErrHandler handles errors
type ErrHandler func(err error)

// WsConfig webservice configuration
type WsConfig struct {
	Endpoint string
	IP       string
	Resolver *net.Resolver
}

func newWsConfig(endpoint string) *WsConfig {
	return &WsConfig{
		Endpoint: endpoint,
	}
}

func (cfg *WsConfig) WithIP(ip string) {
	cfg.IP = ip
}

func (cfg *WsConfig) WithResolver(resolver *net.Resolver) {
	cfg.Resolver = resolver
}

var wsServe = func(cfg *WsConfig, handler WsHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	var Dialer websocket.Dialer
	if cfg.IP == "" {
		Dialer = websocket.Dialer{
			Proxy:             http.ProxyFromEnvironment,
			HandshakeTimeout:  45 * time.Second,
			EnableCompression: false,
		}
	} else {
		Dialer = websocket.Dialer{
			NetDial: func(network, addr string) (net.Conn, error) {
				localAddr, err := net.ResolveTCPAddr("tcp", cfg.IP+":0") // 替换为您的出口IP地址
				if err != nil {
					return nil, err
				}
				var d net.Dialer
				if cfg.Resolver == nil {
					d = net.Dialer{
						LocalAddr: localAddr,
						Resolver:  net.DefaultResolver,
					}
				} else {
					d = net.Dialer{
						LocalAddr: localAddr,
						Resolver:  cfg.Resolver,
					}
				}

				return d.Dial(network, addr)
			},
			HandshakeTimeout:  45 * time.Second,
			EnableCompression: false,
		}
	}

	c, _, err := Dialer.Dial(cfg.Endpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	c.SetReadLimit(655350)
	doneC = make(chan struct{})
	stopC = make(chan struct{})
	go func() {
		// This function will exit either on error from
		// websocket.Conn.ReadMessage or when the stopC channel is
		// closed by the client.
		defer close(doneC)
		if WebsocketKeepalive {
			keepAlive(c, WebsocketTimeout)
		}
		// Wait for the stopC channel to be closed.  We do that in a
		// separate goroutine because ReadMessage is a blocking
		// operation.
		silent := false
		go func() {
			select {
			case <-stopC:
				silent = true
			case <-doneC:
			}
			c.Close()
		}()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				if !silent {
					errHandler(err)
				}
				return
			}
			handler(message)
		}
	}()
	return
}

func keepAlive(c *websocket.Conn, timeout time.Duration) {
	ticker := time.NewTicker(timeout)

	lastResponse := time.Now()
	c.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})

	go func() {
		defer ticker.Stop()
		for {
			deadline := time.Now().Add(10 * time.Second)
			err := c.WriteControl(websocket.PingMessage, []byte{}, deadline)
			if err != nil {
				return
			}
			<-ticker.C
			if time.Since(lastResponse) > timeout {
				c.Close()
				return
			}
		}
	}()
}
