package nserver

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/go-nats"
	"reflect"
	"runtime/debug"
	"strings"
)

type Handler interface{}
type MiddleFunc func(*ReqContext, func())

type NServer struct {
	id       string
	handlers map[string]Handler
	middles  []MiddleFunc
	queue    string
	nc       *nats.Conn
}

func New(queue, id, nats_addr string) *NServer {
	c, err := DialNats(nats_addr)
	if err != nil {
		panic("connect to nats error!!!")
	}

	return &NServer{
		id:       id,
		handlers: map[string]Handler{},
		middles:  make([]MiddleFunc, 0),
		queue:    queue,
		nc:       c,
	}
}

func DialNats(addr string) (*nats.Conn, error) {
	c, err := nats.Connect(addr)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *NServer) Router(pattern string, h Handler) {
	s.handlers[pattern] = h
}

func (s *NServer) Use(funcs ...MiddleFunc) {
	for _, f := range funcs {
		s.middles = append(s.middles, f)
	}
}

func (s *NServer) Serving() {
	for subj, _ := range s.handlers {
		subject := s.queue + "." + subj + ".*"
		q := s.queue + ":" + subj
		s.nc.QueueSubscribe(subject, q, func(msg *nats.Msg) {
			go s.dispatcher(msg)
		})
	}

	select {}
}

type Subject struct {
	App    string
	Module string
	Method string
}

func parseSubject(subj string) *Subject {
	n := strings.Split(subj, ".")
	a, m, act := n[0], n[1], n[2]
	return &Subject{a, m, act}
}

func (s *NServer) dispatcher(msg *nats.Msg) {
	var rc *ReqContext
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			if rc != nil {
				rc.Write([]byte(fmt.Sprintf("server internal error: %v", r)))
			}
		}
	}()

	subj := parseSubject(msg.Subject)
	rc = &ReqContext{
		Module: subj.Module,
		Method: subj.Method,
		reply:  msg.Reply,
		server: s,
	}
	if len(msg.Data) > 0 {
		err := json.Unmarshal(msg.Data, &rc.Param)
		if err != nil {
			rc.Write([]byte(fmt.Sprintf("params should be json format")))
			return
		}
	}

	if req_id := rc.Param.GetString("req_id"); req_id != "" {
		rc.id = req_id
	}

	i := 0
	var next func()
	next = func() {
		if i < len(s.middles) {
			i++
			s.middles[i-1](rc, next)
		} else {
			s.do(rc)
		}
	}
	next()
}

func (s *NServer) do(rc *ReqContext) {
	h, ok := s.handlers[rc.Module]
	if !ok {
		rc.Write([]byte(fmt.Sprintf("module %s not found", rc.Module)))
		return
	}
	action := reflect.ValueOf(h).MethodByName(strings.Title(rc.Method))
	if !action.IsValid() {
		rc.Write([]byte(fmt.Sprintf("method %s not found", rc.Method)))
		return
	}
	args := []reflect.Value{reflect.ValueOf(rc)}
	action.Call(args)
}
