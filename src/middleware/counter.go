package middleware

import (
	"fmt"
	"net/http"
	"time"
)

type Counter struct {
	totalCount uint32;
	todayCount uint32;
	currentHourCount uint32;
	currentMinuteCount uint32;
	lastTime time.Time;
}

func NewCounter() *Counter {
	return &Counter{
		0, 0, 0, 0, time.Now(),
	}
}

func (this *Counter) Handle(req *http.Request, url string) bool {
	now := time.Now()
	if now.Day() != this.lastTime.Day() {
		this.todayCount = 0	
		this.currentHourCount = 0
		this.currentMinuteCount = 0
	} else if now.Hour() != this.lastTime.Hour() {
		this.currentHourCount = 0
		this.currentMinuteCount = 0
	} else if now.Minute() != this.lastTime.Minute() {
		this.currentMinuteCount = 0
	}
	this.lastTime = now

	this.totalCount += 1
	this.todayCount += 1
	this.currentHourCount += 1
	this.currentMinuteCount += 1

	// fmt.Printf("%+v", this)
	return true
}