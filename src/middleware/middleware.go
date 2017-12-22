package middleware

import (
	"net/http"
)

type Middleware interface {
	Handle(req *RequestHolder) bool
}

type HeadMiddleware struct {
	logger Logger;
	counter Counter;
}

type RequestHolder struct {
	Request *http.Request;
	Url string;
	BindId int32;
	GroupId int32;
}

func NewRequestHolder(request *http.Request) RequestHolder {
	return RequestHolder{
		Request: request,
	}
}

func NewHeadMiddleware() *HeadMiddleware {
	return &HeadMiddleware{

	}
}

func Initialize() {
	go DumpLog()
}

func (this *HeadMiddleware) Handle(requestHolder *RequestHolder) bool {
	this.logger.Handle(requestHolder)
	this.counter.Handle(requestHolder)

	return true
}