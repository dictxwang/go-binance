package portfolio

import (
	"context"
	"encoding/json"
	"net/http"
)

type UmSetLeverageService struct {
	c        *Client
	symbol   string
	leverage int64
}

// Symbol set symbol
func (s *UmSetLeverageService) Symbol(symbol string) *UmSetLeverageService {
	s.symbol = symbol
	return s
}

// Leverage set leverage
func (s *UmSetLeverageService) Leverage(leverage int64) *UmSetLeverageService {
	s.leverage = leverage
	return s
}

// Do send request
func (s *UmSetLeverageService) Do(ctx context.Context, opts ...RequestOption) (res *UmSetLeverageResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/papi/v1/um/leverage",
		secType:  secTypeSigned,
	}
	r.setParam("symbol", s.symbol)
	r.setParam("leverage", s.leverage)

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(UmSetLeverageResponse)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type UmSetLeverageResponse struct {
	Symbol           string `json:"symbol"`
	Leverage         int    `json:"leverage"`
	MaxNotionalValue string `json:"maxNotionalValue"`
}
