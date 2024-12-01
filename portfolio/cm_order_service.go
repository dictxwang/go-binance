package portfolio

import (
	"context"
	"encoding/json"
	"net/http"
)

// CmCreateOrderService create order
type CmCreateOrderService struct {
	c                *Client
	symbol           string
	side             SideType
	positionSide     *PositionSideType
	orderType        OrderType
	timeInForce      *TimeInForceType
	quantity         string
	reduceOnly       *bool
	price            *string
	newClientOrderID *string
	newOrderRespType NewOrderRespType
}

// Symbol set symbol
func (s *CmCreateOrderService) Symbol(symbol string) *CmCreateOrderService {
	s.symbol = symbol
	return s
}

// Side set side
func (s *CmCreateOrderService) Side(side SideType) *CmCreateOrderService {
	s.side = side
	return s
}

// PositionSide set side
func (s *CmCreateOrderService) PositionSide(positionSide PositionSideType) *CmCreateOrderService {
	s.positionSide = &positionSide
	return s
}

// Type set type
func (s *CmCreateOrderService) Type(orderType OrderType) *CmCreateOrderService {
	s.orderType = orderType
	return s
}

// TimeInForce set timeInForce
func (s *CmCreateOrderService) TimeInForce(timeInForce TimeInForceType) *CmCreateOrderService {
	s.timeInForce = &timeInForce
	return s
}

// Quantity set quantity
func (s *CmCreateOrderService) Quantity(quantity string) *CmCreateOrderService {
	s.quantity = quantity
	return s
}

// ReduceOnly set reduceOnly
func (s *CmCreateOrderService) ReduceOnly(reduceOnly bool) *CmCreateOrderService {
	s.reduceOnly = &reduceOnly
	return s
}

// Price set price
func (s *CmCreateOrderService) Price(price string) *CmCreateOrderService {
	s.price = &price
	return s
}

// NewClientOrderID set newClientOrderID
func (s *CmCreateOrderService) NewClientOrderID(newClientOrderID string) *CmCreateOrderService {
	s.newClientOrderID = &newClientOrderID
	return s
}

// NewOrderResponseType set newOrderResponseType
func (s *CmCreateOrderService) NewOrderResponseType(newOrderResponseType NewOrderRespType) *CmCreateOrderService {
	s.newOrderRespType = newOrderResponseType
	return s
}

// CmCreateOrderResponse define create order response
type CmCreateOrderResponse struct {
	ClientOrderID     string           `json:"clientOrderId"` //
	CumQty            string           `json:"cumQty"`        //
	CumBase           string           `json:"cumBase"`       //
	ExecutedQuantity  string           `json:"executedQty"`   //
	OrderID           int64            `json:"orderId"`       //
	AvgPrice          string           `json:"avgPrice"`      //
	OrigQuantity      string           `json:"origQty"`       //
	Price             string           `json:"price"`         //
	ReduceOnly        bool             `json:"reduceOnly"`    //
	Side              SideType         `json:"side"`          //
	PositionSide      PositionSideType `json:"positionSide"`  //
	Status            OrderStatusType  `json:"status"`        //
	Symbol            string           `json:"symbol"`        //
	Pair              string           `json:"pair"`
	TimeInForce       TimeInForceType  `json:"timeInForce"`                 //
	Type              OrderType        `json:"type"`                        //
	UpdateTime        int64            `json:"updateTime"`                  // update time
	RateLimitOrder10s string           `json:"rateLimitOrder10s,omitempty"` //
	RateLimitOrder1m  string           `json:"rateLimitOrder1m,omitempty"`  //
}

func (s *CmCreateOrderService) createOrder(ctx context.Context, endpoint string, opts ...RequestOption) (data []byte, header *http.Header, err error) {

	r := &request{
		method:   http.MethodPost,
		endpoint: endpoint,
		secType:  secTypeSigned,
	}
	m := params{
		"symbol":           s.symbol,
		"side":             s.side,
		"type":             s.orderType,
		"newOrderRespType": s.newOrderRespType,
	}
	if s.quantity != "" {
		m["quantity"] = s.quantity
	}
	if s.positionSide != nil {
		m["positionSide"] = *s.positionSide
	}
	if s.timeInForce != nil {
		m["timeInForce"] = *s.timeInForce
	}
	if s.reduceOnly != nil {
		m["reduceOnly"] = *s.reduceOnly
	}
	if s.price != nil {
		m["price"] = *s.price
	}
	if s.newClientOrderID != nil {
		m["newClientOrderId"] = *s.newClientOrderID
	}
	r.setFormParams(m)
	data, header, err = s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []byte{}, &http.Header{}, err
	}
	return data, header, nil
}

// Do send request
func (s *CmCreateOrderService) Do(ctx context.Context, opts ...RequestOption) (res *CmCreateOrderResponse, err error) {
	data, header, err := s.createOrder(ctx, "/papi/v1/cm/order", opts...)
	if err != nil {
		return nil, err
	}
	res = new(CmCreateOrderResponse)
	err = json.Unmarshal(data, res)

	res.RateLimitOrder10s = header.Get("X-Mbx-Order-Count-10s")
	res.RateLimitOrder1m = header.Get("X-Mbx-Order-Count-1m")

	if err != nil {
		return nil, err
	}
	return res, nil
}

// CmOrder define order info
type CmOrder struct {
	AvgPrice         string           `json:"avgPrice"`
	ClientOrderID    string           `json:"clientOrderId"`
	CumQuote         string           `json:"cumQuote"`
	ExecutedQuantity string           `json:"executedQty"`
	OrderID          int64            `json:"orderId"`
	OrigQuantity     string           `json:"origQty"`
	OrigType         OrderType        `json:"origType"`
	Price            string           `json:"price"`
	ReduceOnly       bool             `json:"reduceOnly"`
	Side             SideType         `json:"side"`
	PositionSide     PositionSideType `json:"positionSide"`
	Status           OrderStatusType  `json:"status"`
	Symbol           string           `json:"symbol"`
	Pair             string           `json:"pair"`
	Time             int64            `json:"time"`
	TimeInForce      TimeInForceType  `json:"timeInForce"`
	Type             OrderType        `json:"type"`
	UpdateTime       int64            `json:"updateTime"`
}

// CmListOpenOrdersService list opened orders
type CmListOpenOrdersService struct {
	c      *Client
	symbol string
}

// Symbol set symbol
func (s *CmListOpenOrdersService) Symbol(symbol string) *CmListOpenOrdersService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *CmListOpenOrdersService) Do(ctx context.Context, opts ...RequestOption) (res []*CmOrder, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/cm/openOrders",
		secType:  secTypeSigned,
	}
	if s.symbol != "" {
		r.setParam("symbol", s.symbol)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*CmOrder{}, err
	}
	res = make([]*CmOrder, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*CmOrder{}, err
	}
	return res, nil
}

// CmCancelAllOpenOrdersService cancel all open orders
type CmCancelAllOpenOrdersService struct {
	c      *Client
	symbol string
}

// Symbol set symbol
func (s *CmCancelAllOpenOrdersService) Symbol(symbol string) *CmCancelAllOpenOrdersService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *CmCancelAllOpenOrdersService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	r := &request{
		method:   http.MethodDelete,
		endpoint: "/papi/v1/cm/allOpenOrders",
		secType:  secTypeSigned,
	}
	r.setFormParam("symbol", s.symbol)
	_, _, err = s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return err
	}
	return nil
}

// CmCancelOrderService cancel all open orders
type CmCancelOrderService struct {
	c                 *Client
	symbol            string
	orderID           int64
	origClientOrderID string
}

// CmCancelOrderResponse define response of canceling order
type CmCancelOrderResponse struct {
	AvgPrice         string           `json:"avgPrice"`
	ClientOrderID    string           `json:"clientOrderId"`
	CumQuantity      string           `json:"cumQty"`
	cumBase          string           `json:"cumBase"`
	ExecutedQuantity string           `json:"executedQty"`
	OrderID          int64            `json:"orderId"`
	OrigQuantity     string           `json:"origQty"`
	Price            string           `json:"price"`
	ReduceOnly       bool             `json:"reduceOnly"`
	Side             SideType         `json:"side"`
	PositionSide     PositionSideType `json:"positionSide"`
	Status           OrderStatusType  `json:"status"`
	Symbol           string           `json:"symbol"`
	Pair             string           `json:"pair"`
	TimeInForce      TimeInForceType  `json:"timeInForce"`
	Type             OrderType        `json:"type"`
	UpdateTime       int64            `json:"updateTime"`
}

// Symbol set symbol
func (s *CmCancelOrderService) Symbol(symbol string) *CmCancelOrderService {
	s.symbol = symbol
	return s
}

// OrderID set orderID
func (s *CmCancelOrderService) OrderID(orderID int64) *CmCancelOrderService {
	s.orderID = orderID
	return s
}

// ClientOrderID set clientOrderID
func (s *CmCancelOrderService) ClientOrderID(clientOrderID string) *CmCancelOrderService {
	s.origClientOrderID = clientOrderID
	return s
}

func (s *CmCancelOrderService) Do(ctx context.Context, opts ...RequestOption) (res *CmCancelOrderResponse, err error) {
	r := &request{
		method:   http.MethodDelete,
		endpoint: "/papi/v1/cm/order",
		secType:  secTypeSigned,
	}
	r.setFormParam("symbol", s.symbol)
	if s.orderID != 0 {
		r.setFormParam("orderId", s.orderID)
	}
	if s.origClientOrderID != "" {
		r.setFormParam("origClientOrderId", s.origClientOrderID)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CmCancelOrderResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
