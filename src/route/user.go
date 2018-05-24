package route

import (
	"geekylx.com/CanteenManagementSystemBackend/src/service"
)

// LoginRequest 登录请求参数
type loginRequest struct {
	Account  string
	Password string
}

type loginResponse struct {
	Type  uint8  `json:"type"`
	Token string `json:"token"`
}

type createUserRequest struct {
	Token       string
	Password    string
	AccountType uint8
}

type createUserResponse struct {
	Account string `json:"account"`
}

type deleteUsersRequest struct {
	Token    string
	Accounts []string
}

type deleteUserResponse struct {
	Result map[string]bool `json:"result"`
}

func login(req interface{}) responseWrapper {
	request, ok := req.(loginRequest)
	if !ok {
		return GenerateErrorResponse(PARAM_TYPE_ERROR_CODE, PARAM_TYPE_ERROR_MESSAGE)
	}
	token, accountType, err := service.Login(request.Account, request.Password)
	if err != nil {
		return GenerateErrorResponse(2, err.Error())
	}
	return GenerateSuccessResponse(loginResponse{
		Token: *token,
		Type:  accountType,
	})
}

func createUser(req interface{}) responseWrapper {
	request, ok := req.(createUserRequest)
	if !ok {
		return GenerateErrorResponse(PARAM_TYPE_ERROR_CODE, PARAM_TYPE_ERROR_MESSAGE)
	}
	account, err := service.CreateUser(request.Token, request.Password, request.AccountType)
	if account == nil || err != nil {
		return GenerateErrorResponse(2, err.Error())
	}
	return GenerateSuccessResponse(createUserResponse{
		Account: *account,
	})
}

func deleteUsers(req interface{}) responseWrapper {
	request, ok := req.(deleteUsersRequest)
	if !ok {
		return GenerateErrorResponse(PARAM_TYPE_ERROR_CODE, PARAM_TYPE_ERROR_MESSAGE)
	}
	deletedUsers, err := service.DeleteUsers(request.Token, request.Accounts)
	if deletedUsers == nil || err != nil {
		return GenerateErrorResponse(2, err.Error())
	}
	return GenerateSuccessResponse(deleteUserResponse{
		Result: deletedUsers,
	})
}
