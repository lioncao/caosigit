package tools

import (
	"fmt"
	"os"
	"time"
)

type LogInterface interface {
	LogFile()
	LogError()
	Log()
}

type simpleLog struct {
	MainLogName, ErrorLogName, LogPath, Logname string
	LogDay                                      int
}

func (this *simpleLog) Init(logname, path string) {
	this.LogPath = path
	this.Logname = logname
	this.LogDay = 0
}

func (this *simpleLog) LogError(fmtstr string, a ...interface{}) {
	t := time.Now()
	_, _, day := t.Date()
	if this.LogDay != day {
		NowData := t.Format("2006-01-02")
		this.MainLogName = fmt.Sprintf("%s/%s%s.log", this.LogPath, this.Logname, NowData)
		this.ErrorLogName = fmt.Sprintf("%s/%s%s.err", this.LogPath, this.Logname, NowData)
	}
	LogStr := fmt.Sprintf(fmtstr, a...)
	NowTime := t.Format("15:04:05")
	Str := fmt.Sprintf("%s: %s\n", NowTime, LogStr)
	f, err := os.OpenFile(this.ErrorLogName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
	defer f.Close()
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf(Str)
		_, feer := f.WriteString(Str)
		if feer != nil {
			fmt.Printf("Write log error: %s", feer)
		}
	}
	mf, merr := os.OpenFile(this.MainLogName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
	defer mf.Close()
	if merr != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf(Str)
		_, mfeer := mf.WriteString(Str)
		if mfeer != nil {
			fmt.Printf("Write log error: %s", mfeer)
		}
	}
}

func (this *simpleLog) LogFile() {

}

func (this *simpleLog) Log(fmtstr string, a ...interface{}) {
	LogStr := fmt.Sprintf(fmtstr, a...)
	t := time.Now()
	_, _, day := t.Date()
	if this.LogDay != day {
		NowData := t.Format("2006-01-02")
		this.MainLogName = fmt.Sprintf("%s/%s%s.log", this.LogPath, this.Logname, NowData)
		this.ErrorLogName = fmt.Sprintf("%s/%s%s.err", this.LogPath, this.Logname, NowData)
	}
	NowTime := t.Format("15:04:05")
	Str := fmt.Sprintf("%s: %s\n", NowTime, LogStr)
	f, err := os.OpenFile(this.MainLogName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
	defer f.Close()
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf(Str)
		_, feer := f.WriteString(Str)
		if feer != nil {
			fmt.Printf("Write log error: %s", feer)
		}
	}

}

var slog *simpleLog

func GetLog() *simpleLog {
	if slog == nil {
		slog = new(simpleLog)
	}
	return slog
}
