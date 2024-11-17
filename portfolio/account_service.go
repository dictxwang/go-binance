package portfolio

import (
	"context"
	"encoding/json"
	"net/http"
)

type GetAllBalanceService struct {
	c *Client
}

// Do send request
func (s *GetAllBalanceService) Do(ctx context.Context, opts ...RequestOption) (res []*Balance, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/balance",
		secType:  secTypeSigned,
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = make([]*Balance, 0)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type Balance struct {
	Asset               string `json:"asset"`
	TotalWalletBalance  string `json:"totalWalletBalance"`
	CrossMarginAsset    string `json:"crossMarginAsset"`
	CrossMarginBorrowed string `json:"crossMarginBorrowed"`
	CrossMarginFree     string `json:"crossMarginFree"`
	CrossMarginInterest string `json:"crossMarginInterest"`
	CrossMarginLocked   string `json:"crossMarginLocked"`
	UmWalletBalance     string `json:"umWalletBalance"`
	UmUnrealizedPNL     string `json:"umUnrealizedPNL"`
	CmWalletBalance     string `json:"cmWalletBalance"`
	CmUnrealizedPNL     string `json:"cmUnrealizedPNL"`
	UpdateTime          int64  `json:"updateTime"`
	NegativeBalance     string `json:"negativeBalance"`
}

type GetMaxBorrowableService struct {
	c     *Client
	asset string
}

func (s *GetMaxBorrowableService) Asset(asset string) *GetMaxBorrowableService {
	s.asset = asset
	return s
}

// Do send request
func (s *GetMaxBorrowableService) Do(ctx context.Context, opts ...RequestOption) (res *MaxBorrowable, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/margin/maxBorrowable",
		secType:  secTypeSigned,
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(MaxBorrowable)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type MaxBorrowable struct {
	Amount      float64 `json:"amount"`
	BorrowLimit float64 `json:"borrowLimit"`
}
