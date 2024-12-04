package portfolio

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/dictxwang/go-binance/common"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

// SideType define side type of order
type SideType string

// PositionSideType define position side type of order
type PositionSideType string

// OrderType define order type
type OrderType string

// TimeInForceType define time in force type of order
type TimeInForceType string

// NewOrderRespType define response JSON verbosity
type NewOrderRespType string

// SelfTradePreventionMode define response JSON verbosity
type SelfTradePreventionMode string

// OrderExecutionType define order execution type
type OrderExecutionType string

// OrderStatusType define order status type
type OrderStatusType string

// BusinessUnit define business unit
type BusinessUnit string

// // SymbolType define symbol type
type SymbolType string

type SideEffectType string

//
//// SymbolStatusType define symbol status type
//type SymbolStatusType string
//
//// SymbolFilterType define symbol filter type
//type SymbolFilterType string
//
//// SideEffectType define side effect type for orders
//type SideEffectType string
//
//// WorkingType define working type
//type WorkingType string
//
//// MarginType define margin type
//type MarginType string
//
//// ContractType define contract type
//type ContractType string

// UserDataEventType define user data event type
type UserDataEventType string

// UserDataEventReasonType define reason type for user data event
type UserDataEventReasonType string

//
//// ForceOrderCloseType define reason type for force order
//type ForceOrderCloseType string

type BNBTransferSide string

// Endpoints
const (
	baseApiMainUrl    = "https://papi.binance.com"
	baseApiTestnetUrl = "https://testnet.binancefuture.com"
)

// Global enums
const (
	SideTypeBuy  SideType = "BUY"
	SideTypeSell SideType = "SELL"

	PositionSideTypeBoth  PositionSideType = "BOTH"
	PositionSideTypeLong  PositionSideType = "LONG"
	PositionSideTypeShort PositionSideType = "SHORT"

	OrderTypeLimit       OrderType = "LIMIT"
	OrderTypeMarket      OrderType = "MARKET"
	OrderTypeLiquidation OrderType = "LIQUIDATION"

	TimeInForceTypeGTC TimeInForceType = "GTC" // Good Till Cancel
	TimeInForceTypeIOC TimeInForceType = "IOC" // Immediate or Cancel
	TimeInForceTypeFOK TimeInForceType = "FOK" // Fill or Kill
	TimeInForceTypeGTX TimeInForceType = "GTX" // Good Till Crossing (Post Only)

	NewOrderRespTypeACK    NewOrderRespType = "ACK"
	NewOrderRespTypeRESULT NewOrderRespType = "RESULT"

	SelfTradePreventionModeNONE SelfTradePreventionMode = "NONE"
	SelfTradePreventionModeET   SelfTradePreventionMode = "EXPIRE_TAKER"
	SelfTradePreventionModeEM   SelfTradePreventionMode = "EXPIRE_MAKER"
	SelfTradePreventionModeEB   SelfTradePreventionMode = "EXPIRE_BOTH"

	OrderExecutionTypeNew        OrderExecutionType = "NEW"
	OrderExecutionTypeCanceled   OrderExecutionType = "CANCELED"
	OrderExecutionTypeCalculated OrderExecutionType = "CALCULATED"
	OrderExecutionTypeExpired    OrderExecutionType = "EXPIRED"
	OrderExecutionTypeTrade      OrderExecutionType = "TRADE"

	OrderStatusTypeNew             OrderStatusType = "NEW"
	OrderStatusTypePartiallyFilled OrderStatusType = "PARTIALLY_FILLED"
	OrderStatusTypeFilled          OrderStatusType = "FILLED"
	OrderStatusTypeCanceled        OrderStatusType = "CANCELED"
	OrderStatusTypeExpired         OrderStatusType = "EXPIRED"
	//OrderStatusTypeRejected        OrderStatusType = "REJECTED"
	//OrderStatusTypeNewInsurance    OrderStatusType = "NEW_INSURANCE"
	//OrderStatusTypeNewADL          OrderStatusType = "NEW_ADL"

	UmBusinessUnit BusinessUnit = "UM"
	CmBusinessUnit BusinessUnit = "CM"

	//
	//SymbolTypeFuture SymbolType = "FUTURE"
	//
	//WorkingTypeMarkPrice     WorkingType = "MARK_PRICE"
	//WorkingTypeContractPrice WorkingType = "CONTRACT_PRICE"
	//
	//SymbolStatusTypePreTrading   SymbolStatusType = "PRE_TRADING"
	//SymbolStatusTypeTrading      SymbolStatusType = "TRADING"
	//SymbolStatusTypePostTrading  SymbolStatusType = "POST_TRADING"
	//SymbolStatusTypeEndOfDay     SymbolStatusType = "END_OF_DAY"
	//SymbolStatusTypeHalt         SymbolStatusType = "HALT"
	//SymbolStatusTypeAuctionMatch SymbolStatusType = "AUCTION_MATCH"
	//SymbolStatusTypeBreak        SymbolStatusType = "BREAK"
	//
	//SymbolFilterTypeLotSize          SymbolFilterType = "LOT_SIZE"
	//SymbolFilterTypePrice            SymbolFilterType = "PRICE_FILTER"
	//SymbolFilterTypePercentPrice     SymbolFilterType = "PERCENT_PRICE"
	//SymbolFilterTypeMarketLotSize    SymbolFilterType = "MARKET_LOT_SIZE"
	//SymbolFilterTypeMaxNumOrders     SymbolFilterType = "MAX_NUM_ORDERS"
	//SymbolFilterTypeMaxNumAlgoOrders SymbolFilterType = "MAX_NUM_ALGO_ORDERS"
	//SymbolFilterTypeMinNotional      SymbolFilterType = "MIN_NOTIONAL"
	//
	//SideEffectTypeNoSideEffect SideEffectType = "NO_SIDE_EFFECT"
	//SideEffectTypeMarginBuy    SideEffectType = "MARGIN_BUY"
	//SideEffectTypeAutoRepay    SideEffectType = "AUTO_REPAY"
	//
	//MarginTypeIsolated MarginType = "ISOLATED"
	//MarginTypeCrossed  MarginType = "CROSSED"
	//
	//ContractTypePerpetual ContractType = "PERPETUAL"

	UserDataEventTypeListenKeyExpired UserDataEventType = "listenKeyExpired"
	UserDataEventTypeAccountUpdate    UserDataEventType = "ACCOUNT_UPDATE"
	UserDataEventTypeOrderTradeUpdate UserDataEventType = "ORDER_TRADE_UPDATE"
	//UserDataEventTypeMarginCall          UserDataEventType = "MARGIN_CALL"
	//UserDataEventTypeAccountConfigUpdate UserDataEventType = "ACCOUNT_CONFIG_UPDATE"

	UserDataEventReasonTypeDeposit             UserDataEventReasonType = "DEPOSIT"
	UserDataEventReasonTypeWithdraw            UserDataEventReasonType = "WITHDRAW"
	UserDataEventReasonTypeOrder               UserDataEventReasonType = "ORDER"
	UserDataEventReasonTypeFundingFee          UserDataEventReasonType = "FUNDING_FEE"
	UserDataEventReasonTypeWithdrawReject      UserDataEventReasonType = "WITHDRAW_REJECT"
	UserDataEventReasonTypeAdjustment          UserDataEventReasonType = "ADJUSTMENT"
	UserDataEventReasonTypeInsuranceClear      UserDataEventReasonType = "INSURANCE_CLEAR"
	UserDataEventReasonTypeAdminDeposit        UserDataEventReasonType = "ADMIN_DEPOSIT"
	UserDataEventReasonTypeAdminWithdraw       UserDataEventReasonType = "ADMIN_WITHDRAW"
	UserDataEventReasonTypeMarginTransfer      UserDataEventReasonType = "MARGIN_TRANSFER"
	UserDataEventReasonTypeMarginTypeChange    UserDataEventReasonType = "MARGIN_TYPE_CHANGE"
	UserDataEventReasonTypeAssetTransfer       UserDataEventReasonType = "ASSET_TRANSFER"
	UserDataEventReasonTypeOptionsPremiumFee   UserDataEventReasonType = "OPTIONS_PREMIUM_FEE"
	UserDataEventReasonTypeOptionsSettleProfit UserDataEventReasonType = "OPTIONS_SETTLE_PROFIT"
	//
	//ForceOrderCloseTypeLiquidation ForceOrderCloseType = "LIQUIDATION"
	//ForceOrderCloseTypeADL         ForceOrderCloseType = "ADL"

	timestampKey  = "timestamp"
	signatureKey  = "signature"
	recvWindowKey = "recvWindow"

	SymbolTypeMargin  SymbolType = "MARGIN"
	SymbolTypeFutures SymbolType = "FUTURES"

	SideEffectTypeNoSideEffect SideEffectType = "NO_SIDE_EFFECT"
	SideEffectTypeMarginBuy    SideEffectType = "MARGIN_BUY"
	SideEffectTypeAutoRepay    SideEffectType = "AUTO_REPAY"

	BNBTransferSideToUM   BNBTransferSide = "TO_UM"
	BNBTransferSideFromUM BNBTransferSide = "FROM_UM"
)

func currentTimestamp() int64 {
	return int64(time.Nanosecond) * time.Now().UnixNano() / int64(time.Millisecond)
}

func newJSON(data []byte) (j *simplejson.Json, err error) {
	j, err = simplejson.NewJson(data)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// getApiEndpoint return the base endpoint of the WS according the UseTestnet flag
func getApiEndpoint() string {
	if UseTestnet {
		return baseApiTestnetUrl
	}
	return baseApiMainUrl
}

// NewClient initialize an API client instance with API key and secret key.
// You should always call this function before using this SDK.
// Services will be created by the form client.NewXXXService().
func NewClient(apiKey, secretKey string) *Client {
	return &Client{
		APIKey:     apiKey,
		SecretKey:  secretKey,
		BaseURL:    getApiEndpoint(),
		UserAgent:  "Binance/golang",
		HTTPClient: http.DefaultClient,
		Logger:     log.New(os.Stderr, "Binance-golang ", log.LstdFlags),
	}
}

// NewClientWithIP initialize an API client instance with API key, secret key and local IP.
func NewClientWithIP(apiKey, secretKey, ip string) *Client {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		log.Fatalf("futures NewClientWithIP ip=%s is invalid", ip)
	}

	dialer := &net.Dialer{
		LocalAddr: &net.TCPAddr{
			IP:   parsedIP, // 设置本地出口 IP 地址
			Port: 0,        // 0 表示随机端口
		},
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, addr)
		},
	}

	return &Client{
		APIKey:    apiKey,
		SecretKey: secretKey,
		BaseURL:   getApiEndpoint(),
		UserAgent: "Binance/golang",
		HTTPClient: &http.Client{
			Transport: transport,
		},
		Logger: log.New(os.Stderr, "Binance-golang ", log.LstdFlags),
	}
}

// NewProxiedClient passing a proxy url
func NewProxiedClient(apiKey, secretKey, proxyUrl string) *Client {
	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		log.Fatal(err)
	}
	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &Client{
		APIKey:    apiKey,
		SecretKey: secretKey,
		BaseURL:   getApiEndpoint(),
		UserAgent: "Binance/golang",
		HTTPClient: &http.Client{
			Transport: tr,
		},
		Logger: log.New(os.Stderr, "Binance-golang ", log.LstdFlags),
	}
}

type doFunc func(req *http.Request) (*http.Response, error)

// Client define API client
type Client struct {
	APIKey     string
	SecretKey  string
	BaseURL    string
	UserAgent  string
	HTTPClient *http.Client
	Debug      bool
	Logger     *log.Logger
	TimeOffset int64
	do         doFunc
}

func (c *Client) debug(format string, v ...interface{}) {
	if c.Debug {
		c.Logger.Printf(format, v...)
	}
}

func (c *Client) parseRequest(r *request, opts ...RequestOption) (err error) {
	// set request options from user
	for _, opt := range opts {
		opt(r)
	}
	err = r.validate()
	if err != nil {
		return err
	}

	fullURL := fmt.Sprintf("%s%s", c.BaseURL, r.endpoint)
	if r.recvWindow > 0 {
		r.setParam(recvWindowKey, r.recvWindow)
	}
	if r.secType == secTypeSigned {
		r.setParam(timestampKey, currentTimestamp()-c.TimeOffset)
	}
	queryString := r.query.Encode()
	body := &bytes.Buffer{}
	bodyString := r.form.Encode()
	header := http.Header{}
	if r.header != nil {
		header = r.header.Clone()
	}
	if bodyString != "" {
		header.Set("Content-Type", "application/x-www-form-urlencoded")
		body = bytes.NewBufferString(bodyString)
	}
	if r.secType == secTypeAPIKey || r.secType == secTypeSigned {
		header.Set("X-MBX-APIKEY", c.APIKey)
	}

	if r.secType == secTypeSigned {
		raw := fmt.Sprintf("%s%s", queryString, bodyString)
		mac := hmac.New(sha256.New, []byte(c.SecretKey))
		_, err = mac.Write([]byte(raw))
		if err != nil {
			return err
		}
		v := url.Values{}
		v.Set(signatureKey, fmt.Sprintf("%x", (mac.Sum(nil))))
		if queryString == "" {
			queryString = v.Encode()
		} else {
			queryString = fmt.Sprintf("%s&%s", queryString, v.Encode())
		}
	}
	if queryString != "" {
		fullURL = fmt.Sprintf("%s?%s", fullURL, queryString)
	}
	c.debug("full url: %s, body: %s", fullURL, bodyString)

	r.fullURL = fullURL
	r.header = header
	r.body = body
	return nil
}

func (c *Client) callAPI(ctx context.Context, r *request, opts ...RequestOption) (data []byte, header *http.Header, err error) {
	err = c.parseRequest(r, opts...)
	if err != nil {
		return []byte{}, &http.Header{}, err
	}
	req, err := http.NewRequest(r.method, r.fullURL, r.body)
	if err != nil {
		return []byte{}, &http.Header{}, err
	}
	req = req.WithContext(ctx)
	req.Header = r.header
	c.debug("request: %#v", req)
	f := c.do
	if f == nil {
		f = c.HTTPClient.Do
	}
	res, err := f(req)
	if err != nil {
		return []byte{}, &http.Header{}, err
	}
	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, &http.Header{}, err
	}
	defer func() {
		cerr := res.Body.Close()
		// Only overwrite the retured error if the original error was nil and an
		// error occurred while closing the body.
		if err == nil && cerr != nil {
			err = cerr
		}
	}()
	c.debug("response: %#v", res)
	c.debug("response body: %s", string(data))
	c.debug("response status code: %d", res.StatusCode)

	if res.StatusCode >= http.StatusBadRequest {
		apiErr := new(common.APIError)
		e := json.Unmarshal(data, apiErr)
		if e != nil {
			c.debug("failed to unmarshal json: %s", e)
		}
		return nil, &http.Header{}, apiErr
	}
	return data, &res.Header, nil
}

// SetApiEndpoint set api Endpoint
func (c *Client) SetApiEndpoint(url string) *Client {
	c.BaseURL = url
	return c
}

// NewCmCommissionRateService init cm commission rate service
func (c *Client) NewCmCommissionRateService() *CmCommissionRateService {
	return &CmCommissionRateService{c: c}
}

// NewCmChangePositionModeService init cm position side dual service
func (c *Client) NewCmChangePositionModeService() *CmChangePositionModeService {
	return &CmChangePositionModeService{c: c}
}

// NewCmGetPositionModeService init cm position side dual service
func (c *Client) NewCmGetPositionModeService() *CmGetPositionModeService {
	return &CmGetPositionModeService{c: c}
}

// NewCmSetLeverageService init cm set leverage service
func (c *Client) NewCmSetLeverageService() *CmSetLeverageService {
	return &CmSetLeverageService{c: c}
}

// NewCmCreateOrderService init cancel all um open orders service
func (c *Client) NewCmCreateOrderService() *CmCreateOrderService {
	return &CmCreateOrderService{c: c}
}

// NewCmListOpenOrdersService init list open orders service
func (c *Client) NewCmListOpenOrdersService() *CmListOpenOrdersService {
	return &CmListOpenOrdersService{c: c}
}

// NewCmCancelOrderService init cancel all um open orders service
func (c *Client) NewCmCancelOrderService() *CmCancelOrderService {
	return &CmCancelOrderService{c: c}
}

// NewCmCancelAllOpenOrdersService init cancel all um open orders service
func (c *Client) NewCmCancelAllOpenOrdersService() *CmCancelAllOpenOrdersService {
	return &CmCancelAllOpenOrdersService{c: c}
}

// NewCmGetAccountService init getting account service
func (c *Client) NewCmGetAccountService() *CmGetAccountService {
	return &CmGetAccountService{c: c}
}

// #### um
// NewUmCommissionRateService init um commission rate service
func (c *Client) NewUmCommissionRateService() *UmCommissionRateService {
	return &UmCommissionRateService{c: c}
}

// NewUmGetPositionModeService init um position side dual service
func (c *Client) NewUmGetPositionModeService() *UmGetPositionModeService {
	return &UmGetPositionModeService{c: c}
}

// NewUmChangePositionModeService init um position side dual service
func (c *Client) NewUmChangePositionModeService() *UmChangePositionModeService {
	return &UmChangePositionModeService{c: c}
}

// NewUmSetLeverageService init um set leverage service
func (c *Client) NewUmSetLeverageService() *UmSetLeverageService {
	return &UmSetLeverageService{c: c}
}

// NewUmCreateOrderService init cancel all um open orders service
func (c *Client) NewUmCreateOrderService() *UmCreateOrderService {
	return &UmCreateOrderService{c: c}
}

// NewListOpenOrdersService init list open orders service
func (c *Client) NewUmListOpenOrdersService() *UmListOpenOrdersService {
	return &UmListOpenOrdersService{c: c}
}

// NewUmCancelAllOpenOrdersService init cancel all um open orders service
func (c *Client) NewUmCancelOrderService() *UmCancelOrderService {
	return &UmCancelOrderService{c: c}
}

// NewUmCancelAllOpenOrdersService init cancel all um open orders service
func (c *Client) NewUmCancelAllOpenOrdersService() *UmCancelAllOpenOrdersService {
	return &UmCancelAllOpenOrdersService{c: c}
}

// NewUmGetAccountService init getting account service
func (c *Client) NewUmGetAccountService() *UmGetAccountService {
	return &UmGetAccountService{c: c}
}

// NewStartUserStreamService init starting user stream service
func (c *Client) NewStartUserStreamService() *StartUserStreamService {
	return &StartUserStreamService{c: c}
}

// NewKeepaliveUserStreamService init keep alive user stream service
func (c *Client) NewKeepaliveUserStreamService() *KeepaliveUserStreamService {
	return &KeepaliveUserStreamService{c: c}
}

// NewCloseUserStreamService init closing user stream service
func (c *Client) NewCloseUserStreamService() *CloseUserStreamService {
	return &CloseUserStreamService{c: c}
}

// NewGetAllBalanceService init get balances service
func (c *Client) NewGetAllBalanceService() *GetAllBalanceService {
	return &GetAllBalanceService{c: c}
}

func (c *Client) NewGetMaxBorrowableService() *GetMaxBorrowableService {
	return &GetMaxBorrowableService{c: c}
}

func (c *Client) NewMarginLoanService() *MarginLoanService {
	return &MarginLoanService{c: c}
}

func (c *Client) NewMarginRepayService() *MarginRepayService {
	return &MarginRepayService{c: c}
}

func (c *Client) NewMarginCreateOrderService() *MarginCreateOrderService {
	return &MarginCreateOrderService{c: c}
}

func (c *Client) NewMarginCancelOrderService() *MarginCancelOrderService {
	return &MarginCancelOrderService{c: c}
}

func (c *Client) NewGetAccountService() *GetAccountService {
	return &GetAccountService{c: c}
}

func (c *Client) NewBNBTransferService() *BNBTransferService {
	return &BNBTransferService{c: c}
}

func (c *Client) NewGetRepayFuturesSwitchService() *GetRepayFuturesSwitchService {
	return &GetRepayFuturesSwitchService{c: c}
}

func (c *Client) NewChangeRepayFuturesSwitchService() *ChangeRepayFuturesSwitchService {
	return &ChangeRepayFuturesSwitchService{c: c}
}

func (c *Client) NewAssetCollectionService() *AssetCollectionService {
	return &AssetCollectionService{c: c}
}
