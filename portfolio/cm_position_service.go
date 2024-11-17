package portfolio

import (
	"context"
	"encoding/json"
	"net/http"
)

type CmChangePositionModeService struct {
	c                *Client
	dualSidePosition bool
}

// DualSidePosition set positionSide
func (s *CmChangePositionModeService) DualSidePosition(dualSidePosition bool) *CmChangePositionModeService {
	s.dualSidePosition = dualSidePosition
	return s
}

// Do send request
func (s *CmChangePositionModeService) Do(ctx context.Context, opts ...RequestOption) (res *CmChangePositionMode, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/papi/v1/cm/positionSide/dual",
		secType:  secTypeSigned,
	}
	r.setParam("dualSidePosition", s.dualSidePosition)

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CmChangePositionMode)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type CmChangePositionMode struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type CmGetPositionModeService struct {
	c *Client
}

// Do send request
func (s *CmGetPositionModeService) Do(ctx context.Context, opts ...RequestOption) (res *CmGetPositionMode, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/cm/positionSide/dual",
		secType:  secTypeSigned,
	}

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CmGetPositionMode)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type CmGetPositionMode struct {
	DualSidePosition bool `json:"dualSidePosition"`
}

type CmSetLeverageService struct {
	c        *Client
	symbol   string
	leverage int64
}

// Symbol set symbol
func (s *CmSetLeverageService) Symbol(symbol string) *CmSetLeverageService {
	s.symbol = symbol
	return s
}

// Leverage set leverage
func (s *CmSetLeverageService) Leverage(leverage int64) *CmSetLeverageService {
	s.leverage = leverage
	return s
}

// Do send request
func (s *CmSetLeverageService) Do(ctx context.Context, opts ...RequestOption) (res *CmSetLeverageResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/papi/v1/cm/leverage",
		secType:  secTypeSigned,
	}
	r.setParam("symbol", s.symbol)
	r.setParam("leverage", s.leverage)

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CmSetLeverageResponse)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type CmSetLeverageResponse struct {
	Symbol   string `json:"symbol"`
	Leverage int    `json:"leverage"`
	MaxQty   string `json:"maxQty"`
}
