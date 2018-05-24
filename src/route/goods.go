package route

import (
	"geekylx.com/CanteenManagementSystemBackend/src/service"
)

type addGoodsRequest struct {
	Token   string
	Species uint64
	Price   float64
	Name    string
}

type addGoodsResponse struct {
	Success bool `json:"success"`
}

func addGoods(req interface{}) responseWrapper {
	request, ok := req.(addGoodsRequest)
	if !ok {
		return GenerateErrorResponse(PARAM_TYPE_ERROR_CODE, PARAM_TYPE_ERROR_MESSAGE)
	}
	success, err := service.AddGoods(request.Token, request.Species, request.Price, request.Name)
	if err != nil {
		return GenerateErrorResponse(2, err.Error())
	}
	return GenerateSuccessResponse(addGoodsResponse{
		Success: success,
	})
}
