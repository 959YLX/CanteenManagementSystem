package route

import (
	"geekylx.com/CanteenManagementSystemBackend/src/service"
)

type rechargeRequest struct {
	Token   string
	Account string
	Money   float64
}

type dealResponse struct {
	Success   bool    `json:"success"`
	Remaining float64 `json:"remaining"`
}

type consumeRequest struct {
	Token   string
	Account string
	GoodsID uint
}

type transferAccountRequest struct {
	Token     string
	ToAccount string
	Money     float64
}

type selectFlowingWaterRequest struct {
	Token string
}

func recharge(req interface{}) responseWrapper {
	request, ok := req.(rechargeRequest)
	if !ok {
		return GenerateErrorResponse(PARAM_TYPE_ERROR_CODE, PARAM_TYPE_ERROR_MESSAGE)
	}
	success, remaining, err := service.Recharge(request.Token, request.Account, request.Money)
	if !success || err != nil {
		return GenerateErrorResponse(2, err.Error())
	}
	return GenerateSuccessResponse(dealResponse{
		Success:   success,
		Remaining: remaining,
	})
}

func consume(req interface{}) responseWrapper {
	request, ok := req.(consumeRequest)
	if !ok {
		return GenerateErrorResponse(PARAM_TYPE_ERROR_CODE, PARAM_TYPE_ERROR_MESSAGE)
	}
	success, remaining, err := service.Consume(request.Token, request.Account, request.GoodsID)
	if !success || err != nil {
		return GenerateErrorResponse(2, err.Error())
	}
	return GenerateSuccessResponse(dealResponse{
		Success:   success,
		Remaining: remaining,
	})
}

func transferAccount(req interface{}) responseWrapper {
	request, ok := req.(transferAccountRequest)
	if !ok {
		return GenerateErrorResponse(PARAM_TYPE_ERROR_CODE, PARAM_TYPE_ERROR_MESSAGE)
	}
	success, remaining, err := service.TransferAccount(request.Token, request.ToAccount, request.Money)
	if !success || err != nil {
		return GenerateErrorResponse(2, err.Error())
	}
	return GenerateSuccessResponse(dealResponse{
		Success:   success,
		Remaining: remaining,
	})
}
