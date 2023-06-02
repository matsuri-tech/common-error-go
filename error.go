package merrors

import (
	"encoding/json"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
)

type ErrorType string

const callerSkipCount = 3

type CommonError struct {
	StatusCode int
	Msg        string
	StackTrace string
	ErrorType
}

type ErrorResponse struct {
	Error     string
	ErrorType ErrorType
}

type ErrorResponseCamelCase struct {
	Error     string
	ErrorType ErrorType `json:"errorType"`
}

type ErrorResponseSnakeCase struct {
	Error     string
	ErrorType ErrorType `json:"error_type"`
}

// UnmarshalJSON で上書き
// Go以外の言語についてはKey名がsnake_caseになるものがあるので、そのAPIのレスポンスのハンドリングで困ることがあるから。
func (r *ErrorResponse) UnmarshalJSON(data []byte) error {
	var c ErrorResponseCamelCase
	if err := json.Unmarshal(data, &c); err != nil {
		return err
	}

	if c.ErrorType != "" {
		r.Error = c.Error
		r.ErrorType = c.ErrorType
		return nil
	}

	var s ErrorResponseSnakeCase
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	r.Error = s.Error
	r.ErrorType = s.ErrorType
	return nil
}

func ErrorByStatusCode(statusCode int, msg string, errorType ErrorType) CommonError {
	switch statusCode {
	case http.StatusNotFound:
		return ErrorNotFound(msg, errorType)
	case http.StatusBadRequest:
		return ErrorBadRequest(msg, errorType)
	case http.StatusUnauthorized:
		return ErrorUnauthorized(msg, errorType)
	default:
		return ErrorInternalServerError("internal server error", errorType)
	}
}

func NewCommonError(statusCode int, msg string, errorType ErrorType) CommonError {
	var s = ""
	for i := 2; i >= 0; i-- {
		_, file, line, _ := runtime.Caller(callerSkipCount + i)
		s = s + file + ":" + strconv.Itoa(line) + " "
	}
	return CommonError{
		StatusCode: statusCode,
		Msg:        msg,
		StackTrace: s,
		ErrorType:  errorType,
	}
}

func ErrorNotFound(msg string, errType ErrorType) CommonError {
	return NewCommonError(http.StatusNotFound, msg, errType)
}

func ErrorUnauthorized(msg string, errType ErrorType) CommonError {
	return NewCommonError(http.StatusUnauthorized, msg, errType)
}

func ErrorBadRequest(msg string, errType ErrorType) CommonError {
	return NewCommonError(http.StatusBadRequest, msg, errType)
}

func ErrorInternalServerError(msg string, errType ErrorType) CommonError {
	return NewCommonError(http.StatusInternalServerError, msg, errType)
}

func (e CommonError) Error() string {
	return string(e.ErrorType) + ": " + e.Msg
}

// 内部向けのスタックトレースとかを表示する
func (e CommonError) InternalErrorJson() map[string]interface{} {
	json := map[string]interface{}{}
	json["type"] = string(e.ErrorType)
	json["msg"] = e.Msg
	json["stackTrace"] = e.StackTrace
	return json
}

// 以下の条件のいずれかを満たす場合にはtrue、それ以外の場合にはfalseを返す
// - x, y がDeepEqualである
// - x, y がいずれもCommonErrorであり、かつErrorTypeが同一
// 異常系のテスト時に、期待したErrorTypeが返っているか確認するために使用することを想定
func ErrorTypeEqual(x error, y error) bool {
	if reflect.DeepEqual(x, y) {
		return true
	}

	cErrX, ok := x.(CommonError)
	if !ok {
		return false
	}

	cErrY, ok := y.(CommonError)
	if !ok {
		return false
	}

	if cErrX.ErrorType == cErrY.ErrorType {
		return true
	}

	return false
}
