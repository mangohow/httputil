package httputil

import (
	"net/http"
	"sync"
)

type resResult struct {
	Data    interface{} `json:"data"`
	Status  uint32      `json:"status"`
	Message string      `json:"message"`
}

type Response struct {
	HttpStatus int
	R          resResult
}

var pool = sync.Pool{
	New: func() interface{} {
		return &Response{}
	},
}

// Messager is used to get message from code specified
type Messager interface {
	Message(code uint32) string
}

type MessagerFunc func(uint32) string

func (m MessagerFunc) Message(code uint32) string {
	return m(code)
}

var m Messager

func SetMessager(mg Messager) {
	m = mg
}

func NewResponse(status int, code uint32, data interface{}) *Response {
	if m == nil {
		panic("must set messager")
	}

	response := pool.Get().(*Response)
	response.HttpStatus = status
	response.R.Status = code
	response.R.Message = m.Message(code)
	response.R.Data = data

	return response
}

func putResponse(res *Response) {
	if res != nil {
		res.R.Data = nil
		pool.Put(res)
	}

}

func NewResponseOK(code uint32, data interface{}) *Response {
	return NewResponse(http.StatusOK, code, data)
}

func NewResponseOKND(code uint32) *Response {
	return NewResponse(http.StatusOK, code, nil)
}
