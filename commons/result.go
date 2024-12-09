package commons

import "net/http"

type R struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Ok(data interface{}) *R {
	return &R{Code: http.StatusOK, Msg: "OK", Data: data}
}

func Error(msg string) *R {
	return &R{Code: -1, Msg: msg, Data: nil}
}
