package middleware

import (
	"net/http"
	// "time"
)

type Logger struct {
	
}

func NewLogger() *Logger {
	return &Logger{
	
	}
}

func (this *Logger) Handle(req *http.Request, url string) bool {
	return true
}