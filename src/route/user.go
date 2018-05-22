package route

// LoginRequest 登录请求参数
type loginRequest struct {
	Account  string
	Password string
}

type loginResponse struct {
	Role      int64   `json:"role"`
	Token     string  `json:"tokrn"`
	Remaining float64 `json:"remaining"`
}

type createUserRequest struct {
	Token    string
	Password string
}

type createUserResponse struct {
	Account string `json:"account"`
}

type deleteUsersRequest struct {
	Token   string
	Account []string
}

type deleteUserResponse struct {
	Result map[string]bool `json:"result"`
}

func login(req interface{}) responseWrapper {
	_, ok := req.(loginRequest)
	if !ok {
		return GenerateErrorResponse(PARAM_TYPE_ERROR_CODE, PARAM_TYPE_ERROR_MESSAGE)
	}
	return GenerateSuccessResponse(loginResponse{})
}

func createUser(req interface{}) responseWrapper {
	_, ok := req.(createUserRequest)
	if !ok {
		return GenerateErrorResponse(PARAM_TYPE_ERROR_CODE, PARAM_TYPE_ERROR_MESSAGE)
	}
	return GenerateSuccessResponse(createUserResponse{})
}

func deleteUser(req interface{}) responseWrapper {
	_, ok := req.(deleteUsersRequest)
	if !ok {
		return GenerateErrorResponse(PARAM_TYPE_ERROR_CODE, PARAM_TYPE_ERROR_MESSAGE)
	}
	return GenerateSuccessResponse(deleteUserResponse{})
}
