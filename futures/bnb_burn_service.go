package futures

import (
	"context"
	"net/http"
)

// ToggleFuturesBNBBurnService toggle BNB Burn on futures trade
type ToggleFuturesBNBBurnService struct {
	c       *Client
	feeBurn bool
}

// FeeBurn sets the futures bnb burn parameter.
func (s *ToggleFuturesBNBBurnService) FeeBurn(v bool) *ToggleFuturesBNBBurnService {
	s.feeBurn = v
	return s
}

// Do send request
func (s *ToggleFuturesBNBBurnService) Do(ctx context.Context, opts ...RequestOption) error {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/fapi/v1/feeBurn",
		secType:  secTypeSigned,
	}
	r.setFormParams(params{
		"feeBurn": s.feeBurn,
	})
	_, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return err
	}
	return nil
}
