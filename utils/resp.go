package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H struct {
	Code  int         `json:"code"` // 注意这里的标签
	Data  interface{} `json:"data"`
	Msg   string      `json:"msg"`
	Rows  interface{} `json:"rows"`
	Total interface{} `json:"total"`
}

func Resp(w http.ResponseWriter, code int, data interface{}, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	h := H{
		Code: code,
		Data: data,
		Msg:  msg,
	}
	ret, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(ret)
}

func RespOK(w http.ResponseWriter, data interface{}) {
	Resp(w, 0, data, "请求成功")
}

func RespFail(w http.ResponseWriter, message string) {
	Resp(w, -1, nil, "请求失败")
}

func RespOKList(w http.ResponseWriter, data interface{}, total interface{}) {
	RespList(w, 0, "请求成功", data, total)
}

func RespList(w http.ResponseWriter, code int, msg string, data interface{}, total interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	h := H{
		Code:  code,
		Rows:  data,
		Total: total,
	}
	ret, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(ret)
}
