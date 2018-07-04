package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"geekylx.com/CanteenManagementSystemBackend/src/utils"
)

const (
	PARAM_TYPE_ERROR_CODE    = 1
	PARAM_TYPE_ERROR_MESSAGE = "param type error"
)

type httpServiceHandler struct {
	RequestParam    interface{}
	ExecuteFunction func(interface{}) responseWrapper
}

type responseWrapper struct {
	Code int64       `json:"code"`
	Data interface{} `json:"data"`
}

// StartHTTPServer 启动HTTP服务器
func StartHTTPServer(address string, port int) (err error) {
	handlers := initHandlers()
	listenAddress := fmt.Sprintf("%s:%d", address, port)
	server := http.NewServeMux()
	for url, handler := range handlers {
		server.HandleFunc(url, wrapper(handler.RequestParam, handler.ExecuteFunction))
	}
	return http.ListenAndServe(listenAddress, server)
}

func initHandlers() map[string]httpServiceHandler {
	handlers := make(map[string]httpServiceHandler)
	handlers["/user/login"] = httpServiceHandler{
		RequestParam:    loginRequest{},
		ExecuteFunction: login,
	}
	handlers["/user/create"] = httpServiceHandler{
		RequestParam:    createUserRequest{},
		ExecuteFunction: createUser,
	}
	handlers["/user/delete"] = httpServiceHandler{
		RequestParam:    deleteUsersRequest{},
		ExecuteFunction: deleteUsers,
	}
	handlers["/goods/add"] = httpServiceHandler{
		RequestParam:    addGoodsRequest{},
		ExecuteFunction: addGoods,
	}
	handlers["/deal/recharge"] = httpServiceHandler{
		RequestParam:    rechargeRequest{},
		ExecuteFunction: recharge,
	}
	handlers["/deal/consume"] = httpServiceHandler{
		RequestParam:    consumeRequest{},
		ExecuteFunction: consume,
	}
	handlers["/deal/transfer"] = httpServiceHandler{
		RequestParam:    transferAccountRequest{},
		ExecuteFunction: transferAccount,
	}
	handlers["/goods/list"] = httpServiceHandler{
		RequestParam:    goodsListRequest{},
		ExecuteFunction: goodsList,
	}
	handlers["/user/select"] = httpServiceHandler{
		RequestParam:    selectRecordRequest{},
		ExecuteFunction: selectRecord,
	}
	return handlers
}

func wrapper(req interface{}, exec func(interface{}) responseWrapper) func(w http.ResponseWriter, r *http.Request) {
	typeOfRequest := reflect.TypeOf(req)
	reqFieldCount := typeOfRequest.NumField()
	requestParamsMap := make(map[string]reflect.Type, reqFieldCount)
	for index := 0; index < reqFieldCount; index++ {
		field := typeOfRequest.Field(index)
		requestParamsMap[field.Name] = field.Type
	}
	return func(w http.ResponseWriter, r *http.Request) {
		newStruct := reflect.Indirect(reflect.New(typeOfRequest))
		for paramName, paramType := range requestParamsMap {
			formStringValue := r.FormValue(paramName)
			fmt.Println("str ", paramName, " = ", formStringValue)
			if utils.IsStringEmpty(formStringValue) {
				continue
			}
			if coverValue, err := coverStringValueToKind(formStringValue, paramType); err == nil {
				fmt.Println(paramName, " = ", coverValue)
				newStruct.FieldByName(paramName).Set(reflect.ValueOf(coverValue))
			}
		}
		response := exec(newStruct.Interface())
		fmt.Println(response)
		if jsonBytes, err := json.Marshal(response); err == nil {
			w.Write(jsonBytes)
		}
	}
}

func coverStringValueToKind(strValue string, valueType reflect.Type) (value interface{}, err error) {
	kind := valueType.Kind()
	switch kind {
	case reflect.String:
		value = strValue
	case reflect.Int:
		intValue, parseError := strconv.ParseInt(strValue, 10, 0)
		value = int(intValue)
		err = parseError
	case reflect.Int8:
		int8Value, parseError := strconv.ParseInt(strValue, 10, 8)
		value = int8(int8Value)
		err = parseError
	case reflect.Int16:
		int16Value, parseError := strconv.ParseInt(strValue, 10, 16)
		value = int16(int16Value)
		err = parseError
	case reflect.Int32:
		int32Value, parseError := strconv.ParseInt(strValue, 10, 32)
		value = int32(int32Value)
		err = parseError
	case reflect.Int64:
		int64Value, parseError := strconv.ParseInt(strValue, 10, 64)
		value = int64(int64Value)
		err = parseError
	case reflect.Uint:
		uintValue, parseError := strconv.ParseUint(strValue, 10, 0)
		value = uint(uintValue)
		err = parseError
	case reflect.Uint8:
		uint8Value, parseError := strconv.ParseUint(strValue, 10, 8)
		value = uint8(uint8Value)
		err = parseError
	case reflect.Uint16:
		uint16Value, parseError := strconv.ParseUint(strValue, 10, 16)
		value = uint16(uint16Value)
		err = parseError
	case reflect.Uint32:
		uint32Value, parseError := strconv.ParseUint(strValue, 10, 32)
		value = uint32(uint32Value)
		err = parseError
	case reflect.Uint64:
		uint64Value, parseError := strconv.ParseUint(strValue, 10, 64)
		value = uint64(uint64Value)
		err = parseError
	case reflect.Bool:
		if strValue == "false" || strValue == "0" {
			value = false
		} else {
			value = true
		}
	case reflect.Float32:
		float32Value, parseError := strconv.ParseFloat(strValue, 32)
		value = float32(float32Value)
		err = parseError
	case reflect.Float64:
		float64Value, parseError := strconv.ParseFloat(strValue, 64)
		value = float64(float64Value)
		err = parseError
	case reflect.Slice:
		strArrayValue := strings.Split(strValue, ",")
		elements := reflect.Indirect(reflect.New(valueType))
		elementType := valueType.Elem()
		for _, strArrayElement := range strArrayValue {
			if elementValue, err := coverStringValueToKind(strArrayElement, elementType); err == nil {
				elements.Set(reflect.Append(elements, reflect.ValueOf(elementValue)))
			}
		}
		value = elements.Interface()
	default:
		return nil, errors.New("Type Not Support")
	}
	return
}

// GenerateSuccessResponse 成功执行返回结构
func GenerateSuccessResponse(data interface{}) responseWrapper {
	return responseWrapper{
		Code: 0,
		Data: data,
	}
}

// GenerateErrorResponse 统一错误返回结构
func GenerateErrorResponse(code int64, message string) responseWrapper {
	return responseWrapper{
		Code: code,
		Data: message,
	}
}
