package portfolio

import (
	"context"
	"encoding/json"
	"net/http"
)

type MarginLoanService struct {
	c      *Client
	asset  string
	amount float64
}

func (service *MarginLoanService) Asset(asset string) *MarginLoanService {
	service.asset = asset
	return service
}

func (service *MarginLoanService) Amount(amount float64) *MarginLoanService {
	service.amount = amount
	return service
}

// Do send request
func (s *MarginLoanService) Do(ctx context.Context, opts ...RequestOption) (res *MarginLoanResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/papi/v1/marginLoan",
		secType:  secTypeSigned,
	}
	r.setParam("asset", s.asset)
	r.setParam("amount", s.amount)

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(MarginLoanResponse)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type MarginLoanResponse struct {
	TransactionId int64 `json:"tranId"`
}

type MarginRepayService struct {
	c      *Client
	asset  string
	amount float64
}

func (service *MarginRepayService) Asset(asset string) *MarginRepayService {
	service.asset = asset
	return service
}

func (service *MarginRepayService) Amount(amount float64) *MarginRepayService {
	service.amount = amount
	return service
}

// Do send request
func (s *MarginRepayService) Do(ctx context.Context, opts ...RequestOption) (res *RepayLoanResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/papi/v1/repayLoan",
		secType:  secTypeSigned,
	}
	r.setParam("asset", s.asset)
	r.setParam("amount", s.amount)

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(RepayLoanResponse)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type RepayLoanResponse struct {
	TransactionId int64 `json:"tranId"`
}
