package merrors

import (
	"net/http"
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
