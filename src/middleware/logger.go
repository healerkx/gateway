package middleware

import (
	"fmt"
	"os"
	"time"
)

type Logger struct {
	
}

func NewLogger() *Logger {
	return &Logger{
	
	}
}


type LogFile struct {
	file 		*os.File;
	fileName 	string;
}

func NewLogFile() *LogFile {
	return &LogFile{} 
}

func getFileName(groupId int32, dateTime time.Time) string {
	return fmt.Sprintf("%d_%s.log", groupId, dateTime.Format("2006-01-02"))
}

func (this *LogFile) WriteContent(groupId int32, content string) {
	now := time.Now()
	if this.file != nil {
		// TODO: If expired, switch file to new
		if this.fileName != getFileName(groupId, now) {
			this.file.Close()
			this.file = nil
		}
	}

	if this.file == nil {
		fileName := getFileName(groupId, now)
		var err error
		if this.file, err = os.OpenFile(fileName, os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0666); err == nil {
			this.fileName = fileName
		}
	}

	if this.file != nil {
		this.file.WriteString(content + "\n")
		this.file.Sync()	
	}
}


var gLoggerChannel = make(chan *RequestHolder, 10)

var gLogFileMap = make(map[int32]*LogFile)

func logFileContent(requestHolder *RequestHolder) {
	logFile := gLogFileMap[requestHolder.GroupId]
	if logFile == nil {
		logFile = NewLogFile()
		gLogFileMap[requestHolder.GroupId] = logFile
	}
	logFile.WriteContent(requestHolder.GroupId, requestHolder.Url)
}

func DumpLog() {
	for {
		select {
		case rh := <-gLoggerChannel:
			logFileContent(rh)
			// fmt.Printf("LOG: %s;", msg)
		//default:
		//	time.Sleep(0)
		}
	}
}


func (this *Logger) Handle(requestHolder *RequestHolder) bool {
	// TODO: Send content to a goroutine.
	gLoggerChannel <- requestHolder
	return true
}
