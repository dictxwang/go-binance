package binance

import (
	"context"
	"net/http"
)

// ListCollateralAssetsService Get LTV information and collateral limit of flexible loan's collateral assets
type ListCollateralAssetsService struct {
	c              *Client
	collateralCoin string `json:"collateralCoin"`
}

// CollateralCoin set collateralCoin
func (s *ListCollateralAssetsService) CollateralCoin(collateralCoin string) *ListCollateralAssetsService {
	s.collateralCoin = collateralCoin
	return s
}

// Do send request
func (s *ListCollateralAssetsService) Do(ctx context.Context, opts ...RequestOption) (res *CollateralAssetsResponse, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v2/loan/flexible/collateral/data",
	}
	if s.collateralCoin != "" {
		r.setParam("collateralCoin", s.collateralCoin)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CollateralAssetsResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CollateralAssetsResponse response
type CollateralAssetsResponse struct {
	Total int                `json:"total"`
	Rows  []*CollateralAsset `json:"rows"`
}

type CollateralAsset struct {
	CollateralCoin string  `json:"collateralCoin"`
	InitialLTV     float64 `json:"initialLTV"`
	MarginCallLTV  float64 `json:"marginCallLTV"`
	LiquidationLTV float64 `json:"liquidationLTV"`
	MaxLimit       float64 `json:"maxLimit"`
}

// ListLoanAssetsService Get interest rate and borrow limit of flexible loanable assets.
type ListLoanAssetsService struct {
	c        *Client
	loanCoin string `json:"loanCoin"`
}

// LoanCoin set loanCoin
func (s *ListLoanAssetsService) LoanCoin(loanCoin string) *ListLoanAssetsService {
	s.loanCoin = loanCoin
	return s
}

// Do send request
func (s *ListLoanAssetsService) Do(ctx context.Context, opts ...RequestOption) (res *LoanAssetsResponse, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v2/loan/flexible/loanable/data",
	}
	if s.loanCoin != "" {
		r.setParam("loanCoin", s.loanCoin)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(LoanAssetsResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// LoanAssetsResponse response
type LoanAssetsResponse struct {
	Total int          `json:"total"`
	Rows  []*LoanAsset `json:"rows"`
}

type LoanAsset struct {
	LoanCoin             string  `json:"loanCoin"`
	FlexibleInterestRate float64 `json:"flexibleInterestRate"`
	FlexibleMinLimit     float64 `json:"flexibleMinLimit"`
	FlexibleMaxLimit     float64 `json:"flexibleMaxLimit"`
}

// ListOngoingOrdersService Get Flexible Loan Ongoing Orders
type ListOngoingOrdersService struct {
	c              *Client
	loanCoin       string `json:"loanCoin"`
	collateralCoin string `json:"collateralCoin"`
	current        *int   `json:"current"`
	limit          *int   `json:"limit"`
}

// LoanCoin set loanCoin
func (s *ListOngoingOrdersService) LoanCoin(loanCoin string) *ListOngoingOrdersService {
	s.loanCoin = loanCoin
	return s
}

// CollateralCoin set collateralCoin
func (s *ListOngoingOrdersService) CollateralCoin(collateralCoin string) *ListOngoingOrdersService {
	s.collateralCoin = collateralCoin
	return s
}

// Current set current
func (s *ListOngoingOrdersService) Current(current int) *ListOngoingOrdersService {
	s.current = &current
	return s
}

// Limit set limit
func (s *ListOngoingOrdersService) Limit(limit int) *ListOngoingOrdersService {
	s.limit = &limit
	return s
}

// Do send request
func (s *ListOngoingOrdersService) Do(ctx context.Context, opts ...RequestOption) (res *OngoingOrdersResp, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v2/loan/flexible/ongoing/orders",
		secType:  secTypeSigned,
	}

	if s.loanCoin != "" {
		r.setParam("loanCoin", s.loanCoin)
	}
	if s.collateralCoin != "" {
		r.setParam("collateralCoin", s.collateralCoin)
	}
	if s.current != nil {
		r.setParam("current", *s.current)
	}
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(OngoingOrdersResp)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// OngoingOrdersResp response
type OngoingOrdersResp struct {
	Total int             `json:"total"`
	Rows  []*OngoingOrder `json:"rows"`
}

type OngoingOrder struct {
	LoanCoin         string  `json:"loanCoin"`
	TotalDebt        float64 `json:"totalDebt"`
	CollateralCoin   string  `json:"collateralCoin"`
	CollateralAmount float64 `json:"CollateralAmount"`
	CurrentLTV       float64 `json:"currentLTV"`
}

// FlexibleLoanBorrowService Borrow Flexible Loan
type FlexibleLoanBorrowService struct {
	c                *Client
	loanCoin         string   `json:"loanCoin"`
	loanAmount       *float64 `json:"loanAmount"`
	collateralCoin   string   `json:"collateralCoin"`
	collateralAmount *float64 `json:"collateralAmount"`
}

func (s *FlexibleLoanBorrowService) LoanCoin(loanCoin string) *FlexibleLoanBorrowService {
	s.loanCoin = loanCoin
	return s
}

func (s *FlexibleLoanBorrowService) LoanAmount(loanAmount float64) *FlexibleLoanBorrowService {
	s.loanAmount = &loanAmount
	return s
}

func (s *FlexibleLoanBorrowService) CollateralCoin(collateralCoin string) *FlexibleLoanBorrowService {
	s.collateralCoin = collateralCoin
	return s
}

func (s *FlexibleLoanBorrowService) CollateralAmount(collateralAmount float64) *FlexibleLoanBorrowService {
	s.collateralAmount = &collateralAmount
	return s
}

// Do send request
func (s *FlexibleLoanBorrowService) Do(ctx context.Context, opts ...RequestOption) (res *FlexibleLoanBorrowResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v2/loan/flexible/borrow",
		secType:  secTypeSigned,
	}

	r.setParam("loanCoin", s.loanCoin)
	r.setParam("collateralCoin", s.collateralCoin)

	if s.loanAmount != nil {
		r.setParam("loanAmount", s.loanAmount)
	}

	if s.collateralAmount != nil {
		r.setParam("collateralAmount", s.collateralAmount)
	}
	res = new(FlexibleLoanBorrowResponse)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// FlexibleLoanBorrowResponse define borrow response
type FlexibleLoanBorrowResponse struct {
	LoanCoin         string             `json:"loanCoin"`
	LoanAmount       float64            `json:"loanAmount"`
	CollateralCoin   string             `json:"collateralCoin"`
	CollateralAmount float64            `json:"collateralAmount"`
	Status           FlexibleLoanStatus `json:"status"`
}

// FlexibleLoanRepayService Flexible Loan Repay
type FlexibleLoanRepayService struct {
	c                *Client
	loanCoin         string  `json:"loanCoin"`
	collateralCoin   string  `json:"collateralCoin"`
	repayAmount      float64 `json:"repayAmount"`
	collateralReturn *bool   `json:"collateralReturn"` // Default: TRUE. TRUE: Return extra collateral to spot account; FALSE: Keep extra collateral in the order, and lower LTV.
	fullRepayment    *bool   `json:"fullRepayment"`    // Default: FALSE. TRUE: Full repayment; FALSE: Partial repayment, based on loanAmount
}

func (s *FlexibleLoanRepayService) LoanCoin(loanCoin string) *FlexibleLoanRepayService {
	s.loanCoin = loanCoin
	return s
}

func (s *FlexibleLoanRepayService) CollateralCoin(collateralCoin string) *FlexibleLoanRepayService {
	s.collateralCoin = collateralCoin
	return s
}

func (s *FlexibleLoanRepayService) RepayAmount(repayAmount float64) *FlexibleLoanRepayService {
	s.repayAmount = repayAmount
	return s
}

func (s *FlexibleLoanRepayService) CollateralReturn(collateralReturn bool) *FlexibleLoanRepayService {
	s.collateralReturn = &collateralReturn
	return s
}

func (s *FlexibleLoanRepayService) FullRepayment(fullRepayment bool) *FlexibleLoanRepayService {
	s.fullRepayment = &fullRepayment
	return s
}

// Do send request
func (s *FlexibleLoanRepayService) Do(ctx context.Context, opts ...RequestOption) (res *FlexibleLoanRepayResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v2/loan/flexible/repay",
		secType:  secTypeSigned,
	}

	r.setParam("loanCoin", s.loanCoin)
	r.setParam("collateralCoin", s.collateralCoin)
	r.setParam("repayAmount", s.repayAmount)
	if s.collateralReturn != nil {
		r.setParam("collateralReturn", s.collateralReturn)
	}

	if s.fullRepayment != nil {
		r.setParam("fullRepayment", s.fullRepayment)
	}
	res = new(FlexibleLoanRepayResponse)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// FlexibleLoanRepayResponse define borrow response
type FlexibleLoanRepayResponse struct {
	LoanCoin            string                  `json:"loanCoin"`
	CollateralCoin      string                  `json:"collateralCoin"`
	RemainingDebt       float64                 `json:"remainingDebt"`
	RemainingCollateral float64                 `json:"remainingCollateral"`
	FullRepayment       bool                    `json:"fullRepayment"`
	CurrentLTV          float64                 `json:"currentLTV"`
	RepayStatus         FlexibleLoanRepayStatus `json:"repayStatus"`
}

// FlexibleLoanCollateralRepayService Flexible Loan Collateral Repay
type FlexibleLoanCollateralRepayService struct {
	c              *Client
	loanCoin       string  `json:"loanCoin"`
	collateralCoin string  `json:"collateralCoin"`
	repayAmount    float64 `json:"repayAmount"`
	fullRepayment  *bool   `json:"fullRepayment"` // Default: FALSE. TRUE: Full repayment; FALSE: Partial repayment, based on loanAmount
}

func (s *FlexibleLoanCollateralRepayService) LoanCoin(loanCoin string) *FlexibleLoanCollateralRepayService {
	s.loanCoin = loanCoin
	return s
}

func (s *FlexibleLoanCollateralRepayService) CollateralCoin(collateralCoin string) *FlexibleLoanCollateralRepayService {
	s.collateralCoin = collateralCoin
	return s
}

func (s *FlexibleLoanCollateralRepayService) RepayAmount(repayAmount float64) *FlexibleLoanCollateralRepayService {
	s.repayAmount = repayAmount
	return s
}

func (s *FlexibleLoanCollateralRepayService) FullRepayment(fullRepayment bool) *FlexibleLoanCollateralRepayService {
	s.fullRepayment = &fullRepayment
	return s
}

// Do send request
func (s *FlexibleLoanCollateralRepayService) Do(ctx context.Context, opts ...RequestOption) (res *FlexibleLoanCollateralRepayResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v2/loan/flexible/repay/collateral",
		secType:  secTypeSigned,
	}

	r.setParam("loanCoin", s.loanCoin)
	r.setParam("collateralCoin", s.collateralCoin)
	r.setParam("repayAmount", s.repayAmount)

	if s.fullRepayment != nil {
		r.setParam("fullRepayment", s.fullRepayment)
	}
	res = new(FlexibleLoanCollateralRepayResponse)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// FlexibleLoanCollateralRepayResponse define borrow response
type FlexibleLoanCollateralRepayResponse struct {
	LoanCoin            string                  `json:"loanCoin"`
	CollateralCoin      string                  `json:"collateralCoin"`
	RemainingDebt       float64                 `json:"remainingDebt"`
	RemainingCollateral float64                 `json:"remainingCollateral"`
	FullRepayment       bool                    `json:"fullRepayment"`
	CurrentLTV          float64                 `json:"currentLTV"`
	RepayStatus         FlexibleLoanRepayStatus `json:"repayStatus"`
}

// FlexibleLoanAdjustLTVService Flexible Loan Collateral Repay
type FlexibleLoanAdjustLTVService struct {
	c                *Client
	loanCoin         string             `json:"loanCoin"`
	collateralCoin   string             `json:"collateralCoin"`
	adjustmentAmount float64            `json:"adjustmentAmount"`
	direction        AdjustLTVDirection `json:"fullRepayment"` // "ADDITIONAL", "REDUCED"
}

func (s *FlexibleLoanAdjustLTVService) LoanCoin(loanCoin string) *FlexibleLoanAdjustLTVService {
	s.loanCoin = loanCoin
	return s
}

func (s *FlexibleLoanAdjustLTVService) CollateralCoin(collateralCoin string) *FlexibleLoanAdjustLTVService {
	s.collateralCoin = collateralCoin
	return s
}

func (s *FlexibleLoanAdjustLTVService) AdjustmentAmount(adjustmentAmount float64) *FlexibleLoanAdjustLTVService {
	s.adjustmentAmount = adjustmentAmount
	return s
}

func (s *FlexibleLoanAdjustLTVService) Direction(direction AdjustLTVDirection) *FlexibleLoanAdjustLTVService {
	s.direction = direction
	return s
}

// Do send request
func (s *FlexibleLoanAdjustLTVService) Do(ctx context.Context, opts ...RequestOption) (res *FlexibleLoanAdjustLTVResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v2/loan/flexible/adjust/ltv",
		secType:  secTypeSigned,
	}

	r.setParam("loanCoin", s.loanCoin)
	r.setParam("collateralCoin", s.collateralCoin)
	r.setParam("adjustmentAmount", s.adjustmentAmount)
	r.setParam("direction", s.direction)

	res = new(FlexibleLoanAdjustLTVResponse)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// FlexibleLoanAdjustLTVResponse define borrow response
type FlexibleLoanAdjustLTVResponse struct {
	LoanCoin         string             `json:"loanCoin"`
	CollateralCoin   string             `json:"collateralCoin"`
	Direction        string             `json:"direction"`
	AdjustmentAmount float64            `json:"adjustmentAmount"`
	CurrentLTV       float64            `json:"currentLTV"`
	Status           FlexibleLoanStatus `json:"status"`
}

type (
	FlexibleLoanStatus      string
	FlexibleLoanRepayStatus string
	AdjustLTVDirection      string
)

const (
	FlexibleLoanSucceeds         = FlexibleLoanStatus("Succeeds")
	FlexibleLoanFailed           = FlexibleLoanStatus("Failed")
	FlexibleLoanProcessing       = FlexibleLoanStatus("Processing")
	FlexibleLoanRepayRepaid      = FlexibleLoanRepayStatus("Repaid")
	FlexibleLoanRepayRepaying    = FlexibleLoanRepayStatus("Repaying")
	FlexibleLoanRepayFailed      = FlexibleLoanRepayStatus("Failed")
	AdjustLTVDirectionAdditional = AdjustLTVDirection("ADDITIONAL")
	AdjustLTVDirectionReduced    = AdjustLTVDirection("REDUCED")
)
