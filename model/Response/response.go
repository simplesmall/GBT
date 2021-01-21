package Response

import (
	"net/http"
)

type ResponseStructure struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

//未找到响应体
func (res ResponseStructure) NotFound() (result ResponseStructure) {
	return ResponseStructure{
		http.StatusNotFound,
		"Sorry,Not found",
		"",
	}
}

// 错误响应体
func (res ResponseStructure) FailRes(err interface{}) (result ResponseStructure) {
	return ResponseStructure{
		http.StatusInternalServerError,
		"Sometthing wrong in the server end",
		err,
	}
}

// 错误响应体 + msg
func (res ResponseStructure) FailResWithMsg(msg string, err interface{}) (result ResponseStructure) {
	return ResponseStructure{
		http.StatusInternalServerError,
		msg,
		err,
	}
}

// 错误响应体 + msg  + Code
func (res ResponseStructure) FailResWithMsgCode(code int, msg string, err interface{}) (result ResponseStructure) {
	return ResponseStructure{
		code,
		msg,
		err,
	}
}

//带data 正常响应体
func (res ResponseStructure) OKResult(data interface{}) (result ResponseStructure) {
	result.Code = http.StatusOK
	result.Msg = "Normal......"
	result.Data = data
	return result
}

//带data 正常响应体 + msg
func (res ResponseStructure) OKResultWithMsg(msg string, data interface{}) (result ResponseStructure) {
	result.Code = http.StatusOK
	result.Msg = msg
	result.Data = data
	return result
}
