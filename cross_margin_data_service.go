package binance

import (
	"context"
	"net/http"
)

// CrossMarginDataService fetches the interest and loan limit
type CrossMarginDataService struct {
	c        *Client
	vipLevel int
	coin     string
}

// VipLevel sets the vipLevel parameter.
func (s *CrossMarginDataService) VipLevel(vipLevel int) *CrossMarginDataService {
	s.vipLevel = vipLevel
	return s
}

// Coin sets the coin parameter.
func (s *CrossMarginDataService) Coin(coin string) *CrossMarginDataService {
	s.coin = coin
	return s
}

// Do sends the request.
func (s *CrossMarginDataService) Do(ctx context.Context) (*CrossMarginData, error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/margin/crossMarginData",
		secType:  secTypeSigned,
	}
	if s.vipLevel > 0 {
		r.setParam("vipLevel", string(s.vipLevel))
	}
	if s.coin != "" {
		r.setParam("coin", s.coin)
	}

	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return nil, err
	}
	res := new(CrossMarginData)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CrossMarginData represents a response from CrossMarginDataElement.
type CrossMarginData []CrossMarginDataElement

type CrossMarginDataElement struct {
	VipLevel        int      `json:"vipLevel"`
	Coin            string   `json:"coin"`
	TransferIn      bool     `json:"transferIn"`
	Borrowable      bool     `json:"borrowable"`
	DailyInterest   string   `json:"dailyInterest"`
	YearlyInterest  string   `json:"yearlyInterest"`
	BorrowLimit     string   `json:"borrowLimit"`
	MarginAblePairs []string `json:"marginablePairs"`
}
