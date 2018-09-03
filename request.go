package nserver

import (
	"encoding/json"
	"fmt"
	"strconv"
	"net/http"
)

type ReqContext struct {
	Param   Params
	Header  Params
	Cookie  Params
	Module  string
	Method  string
	Resp    *Response
	reply   string
	server  *NServer
	extLogs []string
	newCookies []*http.Cookie
}

type Response struct {
	Code   int         `json:"code"`
	Errmsg string      `json:"errmsg"`
	Result interface{} `json:"result"`
	Cookie []*http.Cookie `json:"COOKIE,omitempty"`
}

type Params map[string]interface{}

// by default return string
func (p Params) Get(key string) string {
	return p.GetString(key)
}

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

func (rc *ReqContext) AppendExtLog(l string) {
	rc.extLogs = append(rc.extLogs, l)
}

func (rc *ReqContext) Write(b []byte) {
	rc.server.nc.Publish(rc.reply, b)
}

func (rc *ReqContext) SetCookie(c *http.Cookie) {
	rc.newCookies = append(rc.newCookies, c)
}

func (rc *ReqContext) Error(code int, result interface{}, errorMsg string) {
	res := NewResponse(code, result)
	res.Errmsg = errorMsg 
	if len(rc.newCookies) == 0 {
		res.Cookie = rc.newCookies
	}
	rc.Resp = res
	jsonByte, err := json.Marshal(res)
	if err != nil {
	}
	rc.server.nc.Publish(rc.reply, jsonByte)
}

func (rc *ReqContext) Json(result interface{}) {
	res := NewResponse(0, result)
	if len(rc.newCookies) > 0 {
		res.Cookie = rc.newCookies
	}
	rc.Resp = res
	jsonByte, err := json.Marshal(res)
	if err != nil {
	}
	rc.server.nc.Publish(rc.reply, jsonByte)
}
