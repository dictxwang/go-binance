package portfolio

import (
	"context"
	"encoding/json"
	"net/http"
)

type CmPositionSideDualService struct {
	c                *Client
	dualSidePosition string
}

// DualSidePosition set positionSide
func (s *CmPositionSideDualService) DualSidePosition(dualSidePosition string) *CmPositionSideDualService {
	s.dualSidePosition = dualSidePosition
	return s
}

// Do send request
func (s *CmPositionSideDualService) Do(ctx context.Context, opts ...RequestOption) (res *CmPositionSideDualResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/papi/v1/cm/positionSide/dual",
		secType:  secTypeSigned,
	}
	r.setParam("dualSidePosition", s.dualSidePosition)

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CmPositionSideDualResponse)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type CmPositionSideDualResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
