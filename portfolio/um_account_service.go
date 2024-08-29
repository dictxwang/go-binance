package portfolio

import (
	"context"
	"encoding/json"
	"net/http"
)

// UmGetAccountService get account info
type UmGetAccountService struct {
	c *Client
}

// Do send request
func (s *UmGetAccountService) Do(ctx context.Context, opts ...RequestOption) (res *UmAccount, err error) {
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

// UmAccount define account info
type UmAccount struct {
	tradeGroupId int64                `json:"tradeGroupId"`
	Assets       []*UmAccountAsset    `json:"assets"`
	Positions    []*UmAccountPosition `json:"positions"`
}

// UmAccountAsset define account asset
type UmAccountAsset struct {
	Asset                  string `json:"asset"`
	CrossWalletBalance     string `json:"crossWalletBalance"`
	CrossUnPnl             string `json:"crossUnPnl"`
	MaintMargin            string `json:"maintMargin"`
	InitialMargin          string `json:"initialMargin"`
	PositionInitialMargin  string `json:"positionInitialMargin"`
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
	UpdateTime             int64  `json:"updateTime"`
}

// UmAccountPosition define account position
type UmAccountPosition struct {
	Symbol                 string           `json:"symbol"`
	InitialMargin          string           `json:"initialMargin"`
	MaintMargin            string           `json:"maintMargin"`
	UnrealizedProfit       string           `json:"unrealizedProfit"`
	PositionInitialMargin  string           `json:"positionInitialMargin"`
	OpenOrderInitialMargin string           `json:"openOrderInitialMargin"`
	Leverage               string           `json:"leverage"`
	EntryPrice             string           `json:"entryPrice"`
	MaxNotional            string           `json:"maxNotional"`
	BidNotional            string           `json:"bidNotional"`
	AskNotional            string           `json:"askNotional"`
	PositionSide           PositionSideType `json:"positionSide"`
	PositionAmt            string           `json:"positionAmt"`
	UpdateTime             int64            `json:"updateTime"`
	BreakEvenPrice         string           `json:"breakEvenPrice"`
}
