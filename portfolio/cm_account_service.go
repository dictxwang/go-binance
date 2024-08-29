package portfolio

import (
	"context"
	"encoding/json"
	"net/http"
)

// CmGetAccountService get account info
type CmGetAccountService struct {
	c *Client
}

// Do send request
func (s *CmGetAccountService) Do(ctx context.Context, opts ...RequestOption) (res *UmAccount, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/um/account",
		secType:  secTypeSigned,
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(UmAccount)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CmAccount define account info
type CmAccount struct {
	Assets    []*UmAccountAsset    `json:"assets"`
	Positions []*UmAccountPosition `json:"positions"`
}

// CmAccountAsset define account asset
type CmAccountAsset struct {
	Asset                  string `json:"asset"`
	CrossWalletBalance     string `json:"crossWalletBalance"`
	CrossUnPnl             string `json:"crossUnPnl"`
	MaintMargin            string `json:"maintMargin"`
	InitialMargin          string `json:"initialMargin"`
	PositionInitialMargin  string `json:"positionInitialMargin"`
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
	UpdateTime             int64  `json:"updateTime"`
}

// CmAccountPosition define account position
type CmAccountPosition struct {
	Symbol                 string           `json:"symbol"`
	PositionAmt            string           `json:"positionAmt"`
	InitialMargin          string           `json:"initialMargin"`
	MaintMargin            string           `json:"maintMargin"`
	UnrealizedProfit       string           `json:"unrealizedProfit"`
	PositionInitialMargin  string           `json:"positionInitialMargin"`
	OpenOrderInitialMargin string           `json:"openOrderInitialMargin"`
	Leverage               string           `json:"leverage"`
	PositionSide           PositionSideType `json:"positionSide"`
	EntryPrice             string           `json:"entryPrice"`
	MaxQuantity            string           `json:"maxQty"`
	UpdateTime             int64            `json:"updateTime"`
	BreakEvenPrice         string           `json:"breakEvenPrice"`
}
