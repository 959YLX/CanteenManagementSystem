package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// StartHttpServer 启动HTTP服务器
func StartHttpServer(address string, port int) (err error) {
	listenAddress := fmt.Sprintf("%s:%d", address, port)
	server := http.NewServeMux()
	return http.ListenAndServe(listenAddress, server)
}

func wrapper(req interface{}, exec func(interface{}) interface{}) func(w http.ResponseWriter, r *http.Request) {
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
			if coverValue, err := coverStringValueToKind(formStringValue, paramType); err == nil {
				fmt.Println(paramName, " = ", coverValue)
				newStruct.FieldByName(paramName).Set(reflect.ValueOf(coverValue))
			}
		}
		if jsonBytes, err := json.Marshal(exec(newStruct.Interface())); err == nil {
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
