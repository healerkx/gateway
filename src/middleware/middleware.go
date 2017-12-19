package middleware

import (
	"net/http"
	
)

type Middleware interface {
	handle(req *http.Request) bool
}
