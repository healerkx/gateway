package middleware

import (
	
)

type Logger struct {
	
}

func NewLogger() *Logger {
	return &Logger{
	
	}
}

type LogContent struct {
	url string;
}

func NewLogContent(url string) LogContent {
	return LogContent {
		url: url,
	}
}


var gLoggerChannel = make(chan LogContent, 10)

func logFileContent(content LogContent) {

}

func DumpLog() {
	for {
		select {
		case msg := <-gLoggerChannel:
			logFileContent(msg)
			// fmt.Printf("LOG: %s;", msg)
		//default:
		//	time.Sleep(0)
		}
	}
}


func (this *Logger) Handle(requestHolder *RequestHolder) bool {
	// TODO: Send content to a goroutine.
	gLoggerChannel <- NewLogContent(requestHolder.Url)
	return true
}
