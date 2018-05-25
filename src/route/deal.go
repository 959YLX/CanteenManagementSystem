package route

import (
	"geekylx.com/CanteenManagementSystemBackend/src/service"
)

type rechargeRequest struct {
	Token   string
	Account string
	Money   float64
}

type rechargeResponse struct {
	Success   bool    `json:"success"`
	Remaining float64 `json:"remaining"`
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
	return GenerateSuccessResponse(rechargeResponse{
		Success:   success,
		Remaining: remaining,
	})
}
