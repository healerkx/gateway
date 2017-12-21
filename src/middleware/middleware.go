package middleware

import (
	"net/http"
)

type Middleware interface {
	Handle(req *http.Request, url string) bool
}

type HeadMiddleware struct {
	logger Logger;
	counter Counter;
}

func NewHeadMiddleware() *HeadMiddleware {
	return &HeadMiddleware{

	}
}

func (this *HeadMiddleware) Handle(req *http.Request, url string) bool {
	this.counter.Handle(req, url)

	return true
}