package errcode

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code    int      `json:"code"`
	Msg     string   `json:"msg"`
	Details []string `json:"details"`
}

var codes  =map[int]string{}

func NewError(code int,msg string)*Error  {
	if _,ok:=codes[code];ok {
		panic(fmt.Sprintf("错误码%d已经存在，请更换一个",code))
	}
	codes[code]=msg
	return &Error{
		Code:    code,
		Msg:     msg,
		Details: nil,
	}
}

func (e *Error) Error() string  {
	return fmt.Sprintf("错误码:%d,错误信息:%s",e.Code,e.Msg)
}
func (e *Error)WithDetails(details ...string)* Error  {
	e.Details=[]string{}
	for _,d:=range details{
		e.Details=append(e.Details,d)
	}
	return  e
}

func (e *Error) StatusCode() int{
	switch e.Code {
	case Success.Code:
		return http.StatusOK
	case ServerError.Code:
		return http.StatusInternalServerError
	case InvalidParams.Code:
		return http.StatusBadRequest
	case UnauthorizedAuthNotExist.Code:
		fallthrough
	case UnauthorizedTokenError.Code:
		fallthrough
	case UnauthorizedTokenGenerate.Code:
		fallthrough
	case UnauthorizedTokenTimeout.Code:
		fallthrough
	case TooManyRequests.Code:
		return http.StatusTooManyRequests
	}
	return http.StatusInternalServerError
}