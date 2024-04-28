package portfolio

import (
	"context"
	"encoding/json"
	"net/http"
)

type UmPositionSideDualService struct {
	c                *Client
	dualSidePosition string
}

// DualSidePosition set positionSide
func (s *UmPositionSideDualService) DualSidePosition(dualSidePosition string) *UmPositionSideDualService {
	s.dualSidePosition = dualSidePosition
	return s
}

// Do send request
func (s *UmPositionSideDualService) Do(ctx context.Context, opts ...RequestOption) (res *UmPositionSideDualResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/papi/v1/um/positionSide/dual",
		secType:  secTypeSigned,
	}
	r.setParam("dualSidePosition", s.dualSidePosition)

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(UmPositionSideDualResponse)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type UmPositionSideDualResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
