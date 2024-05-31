package portfolio

import (
	"encoding/json"
	"fmt"
	"time"
)

// Endpoints
const (
	baseWsMainUrl    = "wss://fstream.binance.com/pm/ws"
	baseWsTestnetUrl = "wss://fstream.binance.com/pm/ws"
)

var (
	// WebsocketTimeout is an interval for sending ping/pong messages if WebsocketKeepalive is enabled
	WebsocketTimeout = time.Second * 60
	// WebsocketKeepalive enables sending ping/pong messages to check the connection stability
	WebsocketKeepalive = false
	// UseTestnet switch all the WS streams from production to the testnet
	UseTestnet = false
)

// getWsEndpoint return the base endpoint of the WS according the UseTestnet flag
func getWsEndpoint() string {
	if UseTestnet {
		return baseWsTestnetUrl
	}
	return baseWsMainUrl
}

// WsBalance define balance
type WsBalance struct {
	Asset              string `json:"a"`
	Balance            string `json:"wb"`
	CrossWalletBalance string `json:"cw"`
	ChangeBalance      string `json:"bc"`
}

// WsPosition define position
type WsPosition struct {
	Symbol              string           `json:"s"`
	Amount              string           `json:"pa"`
	EntryPrice          string           `json:"ep"`
	AccumulatedRealized string           `json:"cr"`
	UnrealizedPnL       string           `json:"up"`
	Side                PositionSideType `json:"ps"`
	BreakEvenPrice      float64          `json:"bep"`
}

// WsAccountUpdate define account update
type WsAccountUpdate struct {
	Reason    UserDataEventReasonType `json:"m"`
	Balances  []WsBalance             `json:"B"`
	Positions []WsPosition            `json:"P"`
}

// WsOrderTradeUpdate define order trade update
type WsOrderTradeUpdate struct {
	Symbol               string             `json:"s"`   // Symbol
	ClientOrderID        string             `json:"c"`   // Client order ID
	Side                 SideType           `json:"S"`   // Side
	Type                 OrderType          `json:"o"`   // Order type
	TimeInForce          TimeInForceType    `json:"f"`   // Time in force
	OriginalQty          string             `json:"q"`   // Original quantity
	OriginalPrice        string             `json:"p"`   // Original price
	AveragePrice         string             `json:"ap"`  // Average price
	StopPrice            string             `json:"sp"`  // Stop price. Please ignore with TRAILING_STOP_MARKET order
	ExecutionType        OrderExecutionType `json:"x"`   // Execution type
	Status               OrderStatusType    `json:"X"`   // Order status
	ID                   int64              `json:"i"`   // Order ID
	LastFilledQty        string             `json:"l"`   // Order Last Filled Quantity
	AccumulatedFilledQty string             `json:"z"`   // Order Filled Accumulated Quantity
	LastFilledPrice      string             `json:"L"`   // Last Filled Price
	CommissionAsset      string             `json:"N"`   // Commission Asset, will not push if no commission
	Commission           string             `json:"n"`   // Commission, will not push if no commission
	TradeTime            int64              `json:"T"`   // Order Trade Time
	TradeID              int64              `json:"t"`   // Trade ID
	BidsNotional         string             `json:"b"`   // Bids Notional
	AsksNotional         string             `json:"a"`   // Asks Notional
	IsMaker              bool               `json:"m"`   // Is this trade the maker side?
	IsReduceOnly         bool               `json:"R"`   // Is this reduce only
	PositionSide         PositionSideType   `json:"ps"`  // Position Side
	RealizedPnL          string             `json:"rp"`  // Realized Profit of the trade
	StrategyType         string             `json:"st"`  // Strategy type, only pushed with conditional order triggered
	StrategyID           int64              `json:"si"`  // StrategyId,only pushed with conditional order triggered
	STP                  string             `json:"V"`   // STP mode
	GTD                  int64              `json:"gtd"` // TIF GTD order auto cancel time
}

// WsUserDataEvent define user data event
type WsUserDataEvent struct {
	Event            UserDataEventType  `json:"e"`
	BusinessUnit     BusinessUnit       `json:"fs"`
	Time             int64              `json:"E"`
	TransactionTime  int64              `json:"T"`
	AccountAlias     string             `json:"i"`
	AccountUpdate    WsAccountUpdate    `json:"a"`
	OrderTradeUpdate WsOrderTradeUpdate `json:"o"`
}

// WsUserDataHandler handle WsUserDataEvent
type WsUserDataHandler func(event *WsUserDataEvent)

func UmWsUserDataServe(listenKey string, handler WsUserDataHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s", getWsEndpoint(), listenKey)
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		fmt.Println(string(message))
		event := new(WsUserDataEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

func UmWsUserDataServeWithIP(ip string, listenKey string, handler WsUserDataHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s", getWsEndpoint(), listenKey)
	cfg := newWsConfig(endpoint)
	cfg.WithIP(ip)

	wsHandler := func(message []byte) {
		event := new(WsUserDataEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsPackedUserDataEvent define user data event
type WsPackedUserDataEvent struct {
	EventType    string `json:"e"`
	EventContent string `json:"EventContent,omitempty"`
}

// WsPackedUserDataHandler handle WsPackedUserDataEvent
type WsPackedUserDataHandler func(event *WsPackedUserDataEvent)

func UmWsPackedUserDataServe(listenKey string, handler WsPackedUserDataHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s", getWsEndpoint(), listenKey)
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		fmt.Println(string(message))
		event := new(WsPackedUserDataEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			errHandler(err)
			return
		}
		event.EventContent = string(message)
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

func UmWsPackedUserDataServeWithIP(ip string, listenKey string, handler WsPackedUserDataHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s", getWsEndpoint(), listenKey)
	cfg := newWsConfig(endpoint)
	cfg.WithIP(ip)

	wsHandler := func(message []byte) {
		event := new(WsPackedUserDataEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			errHandler(err)
			return
		}
		event.EventContent = string(message)
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}
