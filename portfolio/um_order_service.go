package portfolio

import (
	"context"
	"encoding/json"
	"net/http"
)

// UmCreateOrderService create order
type UmCreateOrderService struct {
	c                       *Client
	symbol                  string
	side                    SideType
	positionSide            *PositionSideType
	orderType               OrderType
	timeInForce             *TimeInForceType
	quantity                string
	reduceOnly              *bool
	price                   *string
	newClientOrderID        *string
	newOrderRespType        NewOrderRespType
	selfTradePreventionMode SelfTradePreventionMode
	goodTillDate            int64
}

// Symbol set symbol
func (s *UmCreateOrderService) Symbol(symbol string) *UmCreateOrderService {
	s.symbol = symbol
	return s
}

// Side set side
func (s *UmCreateOrderService) Side(side SideType) *UmCreateOrderService {
	s.side = side
	return s
}

// PositionSide set side
func (s *UmCreateOrderService) PositionSide(positionSide PositionSideType) *UmCreateOrderService {
	s.positionSide = &positionSide
	return s
}

// Type set type
func (s *UmCreateOrderService) Type(orderType OrderType) *UmCreateOrderService {
	s.orderType = orderType
	return s
}

// TimeInForce set timeInForce
func (s *UmCreateOrderService) TimeInForce(timeInForce TimeInForceType) *UmCreateOrderService {
	s.timeInForce = &timeInForce
	return s
}

// Quantity set quantity
func (s *UmCreateOrderService) Quantity(quantity string) *UmCreateOrderService {
	s.quantity = quantity
	return s
}

// ReduceOnly set reduceOnly
func (s *UmCreateOrderService) ReduceOnly(reduceOnly bool) *UmCreateOrderService {
	s.reduceOnly = &reduceOnly
	return s
}

// Price set price
func (s *UmCreateOrderService) Price(price string) *UmCreateOrderService {
	s.price = &price
	return s
}

// NewClientOrderID set newClientOrderID
func (s *UmCreateOrderService) NewClientOrderID(newClientOrderID string) *UmCreateOrderService {
	s.newClientOrderID = &newClientOrderID
	return s
}

// NewOrderResponseType set newOrderResponseType
func (s *UmCreateOrderService) NewOrderResponseType(newOrderResponseType NewOrderRespType) *UmCreateOrderService {
	s.newOrderRespType = newOrderResponseType
	return s
}

// SelfTradePrevention set selfTradePreventionMode
func (s *UmCreateOrderService) SelfTradePrevention(selfTradePreventionMode SelfTradePreventionMode) *UmCreateOrderService {
	s.selfTradePreventionMode = selfTradePreventionMode
	return s
}

// CreateOrderResponse define create order response
type UmCreateOrderResponse struct {
	ClientOrderID           string                  `json:"clientOrderId"`           //
	CumQty                  string                  `json:"cumQty"`                  //
	CumQuote                string                  `json:"cumQuote"`                //
	ExecutedQuantity        string                  `json:"executedQty"`             //
	OrderID                 int64                   `json:"orderId"`                 //
	AvgPrice                string                  `json:"avgPrice"`                //
	OrigQuantity            string                  `json:"origQty"`                 //
	Price                   string                  `json:"price"`                   //
	ReduceOnly              bool                    `json:"reduceOnly"`              //
	Side                    SideType                `json:"side"`                    //
	PositionSide            PositionSideType        `json:"positionSide"`            //
	Status                  OrderStatusType         `json:"status"`                  //
	Symbol                  string                  `json:"symbol"`                  //
	TimeInForce             TimeInForceType         `json:"timeInForce"`             //
	Type                    OrderType               `json:"type"`                    //
	SelfTradePreventionMode SelfTradePreventionMode `json:"selfTradePreventionMode"` // self trading preventation mode
	GoodTillDate            int64                   `json:"goodTillDate"`            // order pre-set auto cancel time for TIF GTD order
	UpdateTime              int64                   `json:"updateTime"`              // update time
}

func (s *UmCreateOrderService) createOrder(ctx context.Context, endpoint string, opts ...RequestOption) (data []byte, header *http.Header, err error) {

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
	if s.selfTradePreventionMode != "" {
		m["selfTradePreventionMode"] = s.selfTradePreventionMode
	}
	if s.goodTillDate != 0 {
		m["goodTillDate"] = s.goodTillDate
	}
	r.setFormParams(m)
	data, header, err = s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []byte{}, &http.Header{}, err
	}
	return data, header, nil
}

// Do send request
func (s *UmCreateOrderService) Do(ctx context.Context, opts ...RequestOption) (res *UmCreateOrderResponse, err error) {
	data, _, err := s.createOrder(ctx, "/papi/v1/um/order", opts...)
	if err != nil {
		return nil, err
	}
	res = new(UmCreateOrderResponse)
	err = json.Unmarshal(data, res)

	if err != nil {
		return nil, err
	}
	return res, nil
}

// UmOrder define order info
type UmOrder struct {
	AvgPrice                string                  `json:"avgPrice"`
	ClientOrderID           string                  `json:"clientOrderId"`
	CumQuote                string                  `json:"cumQuote"`
	ExecutedQuantity        string                  `json:"executedQty"`
	OrderID                 int64                   `json:"orderId"`
	OrigQuantity            string                  `json:"origQty"`
	OrigType                OrderType               `json:"origType"`
	Price                   string                  `json:"price"`
	ReduceOnly              bool                    `json:"reduceOnly"`
	Side                    SideType                `json:"side"`
	PositionSide            PositionSideType        `json:"positionSide"`
	Status                  OrderStatusType         `json:"status"`
	Symbol                  string                  `json:"symbol"`
	Time                    int64                   `json:"time"`
	TimeInForce             TimeInForceType         `json:"timeInForce"`
	Type                    OrderType               `json:"type"`
	UpdateTime              int64                   `json:"updateTime"`
	SelfTradePreventionMode SelfTradePreventionMode `json:"selfTradePreventionMode"`
	GoodTillDate            int64                   `json:"goodTillDate"`
}

// UmListOpenOrdersService list opened orders
type UmListOpenOrdersService struct {
	c      *Client
	symbol string
}

// Symbol set symbol
func (s *UmListOpenOrdersService) Symbol(symbol string) *UmListOpenOrdersService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *UmListOpenOrdersService) Do(ctx context.Context, opts ...RequestOption) (res []*UmOrder, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/um/openOrders",
		secType:  secTypeSigned,
	}
	if s.symbol != "" {
		r.setParam("symbol", s.symbol)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*UmOrder{}, err
	}
	res = make([]*UmOrder, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*UmOrder{}, err
	}
	return res, nil
}

// UmCancelAllOpenOrdersService cancel all open orders
type UmCancelAllOpenOrdersService struct {
	c      *Client
	symbol string
}

// Symbol set symbol
func (s *UmCancelAllOpenOrdersService) Symbol(symbol string) *UmCancelAllOpenOrdersService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *UmCancelAllOpenOrdersService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	r := &request{
		method:   http.MethodDelete,
		endpoint: "/papi/v1/um/allOpenOrders",
		secType:  secTypeSigned,
	}
	r.setFormParam("symbol", s.symbol)
	_, _, err = s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return err
	}
	return nil
}

// UmCancelOrderService cancel all open orders
type UmCancelOrderService struct {
	c                 *Client
	symbol            string
	orderID           int64
	origClientOrderID string
}

// UmCancelOrderResponse define response of canceling order
type UmCancelOrderResponse struct {
	AvgPrice                string                  `json:"avgPrice"`
	ClientOrderID           string                  `json:"clientOrderId"`
	CumQuantity             string                  `json:"cumQty"`
	CumQuote                string                  `json:"cumQuote"`
	ExecutedQuantity        string                  `json:"executedQty"`
	OrderID                 int64                   `json:"orderId"`
	OrigQuantity            string                  `json:"origQty"`
	Price                   string                  `json:"price"`
	ReduceOnly              bool                    `json:"reduceOnly"`
	Side                    SideType                `json:"side"`
	PositionSide            PositionSideType        `json:"positionSide"`
	Status                  OrderStatusType         `json:"status"`
	Symbol                  string                  `json:"symbol"`
	TimeInForce             TimeInForceType         `json:"timeInForce"`
	Type                    OrderType               `json:"type"`
	UpdateTime              int64                   `json:"updateTime"`
	SelfTradePreventionMode SelfTradePreventionMode `json:"selfTradePreventionMode"`
	GoodTillDate            int64                   `json:"goodTillDate"`
}

// Symbol set symbol
func (s *UmCancelOrderService) Symbol(symbol string) *UmCancelOrderService {
	s.symbol = symbol
	return s
}

// OrderID set orderID
func (s *UmCancelOrderService) OrderID(orderID int64) *UmCancelOrderService {
	s.orderID = orderID
	return s
}

// ClientOrderID set clientOrderID
func (s *UmCancelOrderService) ClientOrderID(clientOrderID string) *UmCancelOrderService {
	s.origClientOrderID = clientOrderID
	return s
}

func (s *UmCancelOrderService) Do(ctx context.Context, opts ...RequestOption) (res *UmCancelOrderResponse, err error) {
	r := &request{
		method:   http.MethodDelete,
		endpoint: "/papi/v1/um/order",
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
	res = new(UmCancelOrderResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
