package nserver

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type ReqContext struct {
	Param   Params
	Module  string
	Method  string
	Resp    *Response
	reply   string
	server  *NServer
	id      string
	extLogs []string
}

type Response struct {
	Code   int         `json:"code"`
	Errmsg string      `json:"errmsg"`
	Result interface{} `json:"result"`
}

type Params map[string]interface{}

func (p Params) GetString(key string) string {
	if v, ok := p[key]; ok {
		return v.(string)
	}
	return ""
}

func (p Params) GetInt(key string) (int64, error) {
	if v, ok := p[key]; ok {
		return strconv.ParseInt(v.(string), 10, 64)
	}
	return 0, fmt.Errorf("key:%s not found", key)
}

func NewResponse(code int, result interface{}) *Response {
	res := new(Response)
	res.Code = code
	res.Result = result
	return res
}

func (rc *ReqContext) ExtLogs() []string {
	return rc.extLogs
}

func (rc *ReqContext) ServerID() string {
	return rc.server.id
}

func (rc *ReqContext) App() string {
	return rc.server.queue
}

func (rc *ReqContext) ID() string {
	return rc.id
}

func (rc *ReqContext) AppendExtLog(l string) {
	rc.extLogs = append(rc.extLogs, l)
}

func (rc *ReqContext) Write(b []byte) {
	rc.server.nc.Publish(rc.reply, b)
}

func (rc *ReqContext) Error(code int, result interface{}, errParams map[string]interface{}) {
	res := NewResponse(code, result)
	res.Errmsg = fmt.Sprintf("%+v", errParams)
	rc.Resp = res
	jsonByte, err := json.Marshal(res)
	if err != nil {
	}
	rc.server.nc.Publish(rc.reply, jsonByte)
}

func (rc *ReqContext) Json(result interface{}) {
	res := NewResponse(0, result)
	rc.Resp = res
	jsonByte, err := json.Marshal(res)
	if err != nil {
	}
	rc.server.nc.Publish(rc.reply, jsonByte)
}
