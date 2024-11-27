package portfolio

import (
	"context"
	"encoding/json"
	"net/http"
)

// MarginCreateOrderService create order
type MarginCreateOrderService struct {
	c                       *Client
	symbol                  string
	side                    SideType
	orderType               OrderType
	sideEffectType          *SideEffectType
	timeInForce             *TimeInForceType
	quantity                float64
	price                   float64
	newClientOrderID        string
	newOrderRespType        NewOrderRespType
	selfTradePreventionMode SelfTradePreventionMode
}

func (s *MarginCreateOrderService) Symbol(symbol string) *MarginCreateOrderService {
	s.symbol = symbol
	return s
}

func (s *MarginCreateOrderService) Side(side SideType) *MarginCreateOrderService {
	s.side = side
	return s
}

func (s *MarginCreateOrderService) OrderType(orderType OrderType) *MarginCreateOrderService {
	s.orderType = orderType
	return s
}

func (s *MarginCreateOrderService) SideEffectType(sideEffectType SideEffectType) *MarginCreateOrderService {
	s.sideEffectType = &sideEffectType
	return s
}

func (s *MarginCreateOrderService) TimeInForceType(timeInForceType TimeInForceType) *MarginCreateOrderService {
	s.timeInForce = &timeInForceType
	return s
}

func (s *MarginCreateOrderService) Quantity(quantity float64) *MarginCreateOrderService {
	s.quantity = quantity
	return s
}

func (s *MarginCreateOrderService) Price(price float64) *MarginCreateOrderService {
	s.price = price
	return s
}

func (s *MarginCreateOrderService) NewClientOrderID(newClientOrderID string) *MarginCreateOrderService {
	s.newClientOrderID = newClientOrderID
	return s
}

func (s *MarginCreateOrderService) NewOrderRespType(newOrderRespType NewOrderRespType) *MarginCreateOrderService {
	s.newOrderRespType = newOrderRespType
	return s
}

func (s *MarginCreateOrderService) SelfTradePreventionMode(selfTradePreventionMode SelfTradePreventionMode) *MarginCreateOrderService {
	s.selfTradePreventionMode = selfTradePreventionMode
	return s
}

type MarginCreateOrderResponse struct {
	Symbol                  string                  `json:"symbol"`        //
	OrderID                 int64                   `json:"orderId"`       //
	ClientOrderID           string                  `json:"clientOrderId"` //
	TransactTime            int64                   `json:"transactTime"`
	Price                   string                  `json:"price"`                           //
	SelfTradePreventionMode SelfTradePreventionMode `json:"selfTradePreventionMode"`         // self trading preventation mode
	OrigQuantity            string                  `json:"origQty"`                         //
	ExecutedQuantity        string                  `json:"executedQty"`                     //
	CummulativeQuoteQty     string                  `json:"cummulativeQuoteQty"`             //
	Status                  OrderStatusType         `json:"status"`                          //
	TimeInForce             TimeInForceType         `json:"timeInForce"`                     //
	Type                    OrderType               `json:"type"`                            //
	Side                    SideType                `json:"side"`                            //
	MarginBuyBorrowAsset    string                  `json:"marginBuyBorrowAsset,omitempty"`  //
	MarginBuyBorrowAmount   string                  `json:"marginBuyBorrowAmount,omitempty"` //
}

func (s *MarginCreateOrderService) createOrder(ctx context.Context, endpoint string, opts ...RequestOption) (data []byte, header *http.Header, err error) {

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
	if s.quantity != 0 {
		m["quantity"] = s.quantity
	}
	if s.timeInForce != nil {
		m["timeInForce"] = *s.timeInForce
	}
	if s.price != 0 {
		m["price"] = s.price
	}
	if s.newClientOrderID != "" {
		m["newClientOrderId"] = s.newClientOrderID
	}
	if s.selfTradePreventionMode != "" {
		m["selfTradePreventionMode"] = s.selfTradePreventionMode
	}
	if s.sideEffectType != nil {
		m["sideEffectType"] = *s.sideEffectType
	}

	r.setFormParams(m)
	data, header, err = s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []byte{}, &http.Header{}, err
	}
	return data, header, nil
}

// Do send request
func (s *MarginCreateOrderService) Do(ctx context.Context, opts ...RequestOption) (res *MarginCreateOrderResponse, err error) {
	data, _, err := s.createOrder(ctx, "/papi/v1/margin/order", opts...)
	if err != nil {
		return nil, err
	}
	res = new(MarginCreateOrderResponse)
	err = json.Unmarshal(data, res)

	if err != nil {
		return nil, err
	}
	return res, nil
}

// MarginCancelOrderService cancel order
type MarginCancelOrderService struct {
	c                 *Client
	symbol            string
	orderID           int64
	origClientOrderID string
}

func (s *MarginCancelOrderService) Symbol(symbol string) *MarginCancelOrderService {
	s.symbol = symbol
	return s
}

func (s *MarginCancelOrderService) OrderID(orderID int64) *MarginCancelOrderService {
	s.orderID = orderID
	return s
}

func (s *MarginCancelOrderService) OrigClientOrderID(origClientOrderID string) *MarginCancelOrderService {
	s.origClientOrderID = origClientOrderID
	return s
}

type MarginCancelOrderResponse struct {
	Symbol                  string                  `json:"symbol"`
	OrderID                 int64                   `json:"orderId"`
	OrigClientOrderId       string                  `json:"origClientOrderId"`
	ClientOrderID           string                  `json:"clientOrderId"`
	OrigQuantity            string                  `json:"origQty"`
	Price                   string                  `json:"price"`
	ExecutedQuantity        string                  `json:"executedQty"`
	CummulativeQuoteQty     string                  `json:"cummulativeQuoteQty"`
	Status                  OrderStatusType         `json:"status"`
	TimeInForce             TimeInForceType         `json:"timeInForce"`
	Type                    OrderType               `json:"type"`
	Side                    SideType                `json:"side"`
	SelfTradePreventionMode SelfTradePreventionMode `json:"selfTradePreventionMode"`
}

func (s *MarginCancelOrderService) Do(ctx context.Context, opts ...RequestOption) (res *MarginCancelOrderResponse, err error) {
	r := &request{
		method:   http.MethodDelete,
		endpoint: "/papi/v1/margin/order",
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
	res = new(MarginCancelOrderResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
