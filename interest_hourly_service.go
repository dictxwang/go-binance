package binance

import (
	"context"
	"net/http"
	"strings"
)

// InterestHourlyService fetches the hourly interest
type InterestHourlyService struct {
	c          *Client
	assets     *[]string
	isIsolated bool
}

// Assets sets the assets parameter.
func (s *InterestHourlyService) Assets(assets []string) *InterestHourlyService {
	s.assets = &assets
	return s
}

// IsIsolated sets the isIsolated parameter.
func (s *InterestHourlyService) IsIsolated(isIsolated bool) *InterestHourlyService {
	s.isIsolated = isIsolated
	return s
}

// Do sends the request.
func (s *InterestHourlyService) Do(ctx context.Context) (*HourlyInterest, error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/margin/next-hourly-interest-rate",
		secType:  secTypeSigned,
	}
	r.setParam("assets", strings.Join(*s.assets, ","))
	if s.isIsolated {
		r.setParam("isIsolated", "TRUE")
	} else {
		r.setParam("isIsolated", "FALSE")
	}

	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return nil, err
	}
	res := new(HourlyInterest)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// HourlyInterest represents a response from HourlyInterestElement.
type HourlyInterest []HourlyInterestElement

type HourlyInterestElement struct {
	Asset                  string `json:"asset"`
	NextHourlyInterestRate string `json:"nextHourlyInterestRate"`
}
