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
	err = json.Unmarshal(data, &res)
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
	r.setParam("asset", s.asset)
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
	Amount      string `json:"amount"`
	BorrowLimit string `json:"borrowLimit"`
}

type GetAccountService struct {
	c *Client
}

// Do send request
func (s *GetAccountService) Do(ctx context.Context, opts ...RequestOption) (res *AccountInfo, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/account",
		secType:  secTypeSigned,
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(AccountInfo)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type AccountInfo struct {
	UniMMR                   string `json:"uniMMR"`
	AccountEquity            string `json:"accountEquity"`
	ActualEquity             string `json:"actualEquity"`
	AccountInitialMargin     string `json:"accountInitialMargin"`
	AccountMaintMargin       string `json:"accountMaintMargin"`
	AccountStatus            string `json:"accountStatus"`
	VirtualMaxWithdrawAmount string `json:"virtualMaxWithdrawAmount"`
	TotalAvailableBalance    string `json:"totalAvailableBalance"`
	TotalMarginOpenLoss      string `json:"totalMarginOpenLoss"`
	UpdateTime               int64  `json:"updateTime"`
}

type BNBTransferService struct {
	c            *Client
	amount       float64
	transferSide BNBTransferSide
}

func (c *BNBTransferService) Amount(amount float64) *BNBTransferService {
	c.amount = amount
	return c
}

func (c *BNBTransferService) TransferSide(side BNBTransferSide) *BNBTransferService {
	c.transferSide = side
	return c
}

type BNBTransfer struct {
	TransferId int64 `json:"tranId"`
}

// Do send request
func (s *BNBTransferService) Do(ctx context.Context, opts ...RequestOption) (res *BNBTransfer, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/papi/v1/bnb-transfer",
		secType:  secTypeSigned,
	}
	r.setParam("amount", s.amount)
	r.setParam("transferSide", s.transferSide)

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(BNBTransfer)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type ChangeRepayFuturesSwitchService struct {
	c         *Client
	autoRepay bool
}

func (c *ChangeRepayFuturesSwitchService) AutoRepay(auto bool) *ChangeRepayFuturesSwitchService {
	c.autoRepay = auto
	return c
}

type ChangeRepayFuturesResponse struct {
	Message string `json:"msg"`
}

// Do send request
func (s *ChangeRepayFuturesSwitchService) Do(ctx context.Context, opts ...RequestOption) (res *ChangeRepayFuturesResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/papi/v1/repay-futures-switch",
		secType:  secTypeSigned,
	}
	r.setParam("autoRepay", s.autoRepay)

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(ChangeRepayFuturesResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type GetRepayFuturesSwitchService struct {
	c *Client
}

type GetRepayFuturesResponse struct {
	AutoRepay bool `json:"autoRepay"`
}

// Do send request
func (s *GetRepayFuturesSwitchService) Do(ctx context.Context, opts ...RequestOption) (res *GetRepayFuturesResponse, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/repay-futures-switch",
		secType:  secTypeSigned,
	}

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(GetRepayFuturesResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type AssetCollectionService struct {
	c     *Client
	asset string
}

func (s *AssetCollectionService) Asset(asset string) *AssetCollectionService {
	s.asset = asset
	return s
}

type AssetCollectionResponse struct {
	Message string `json:"msg"`
}

// Do send request
func (s *AssetCollectionService) Do(ctx context.Context, opts ...RequestOption) (res *AssetCollectionResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/papi/v1/asset-collection",
		secType:  secTypeSigned,
	}
	r.setParam("asset", s.asset)

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(AssetCollectionResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
