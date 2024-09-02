package delivery

import (
	"context"
	"encoding/json"
	"net/http"
)

type CommissionRateService struct {
	c      *Client
	symbol string
}

// Symbol set symbol
func (service *CommissionRateService) Symbol(symbol string) *CommissionRateService {
	service.symbol = symbol
	return service
}

// Do send request
func (service *CommissionRateService) Do(ctx context.Context, opts ...RequestOption) (res *CommissionRate, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/dapi/v1/commissionRate",
		secType:  secTypeSigned,
	}
	if service.symbol != "" {
		r.setParam("symbol", service.symbol)
	}
	data, err := service.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CommissionRate)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type CommissionRate struct {
	Symbol              string `json:"symbol"`
	MakerCommissionRate string `json:"makerCommissionRate"`
	TakerCommissionRate string `json:"takerCommissionRate"`
}
