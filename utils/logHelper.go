package utils

import (
	"golang.org/x/sys/windows/svc/eventlog"
)

// 不依赖其它文件的调试，主要用来调试服务
var IsDebug bool = false

type MyLog struct {
}

func (myLog *MyLog) Info(eid uint32, msg string) error {
	if !IsDebug {
		return nil
	}
	return elog.Info(eid, msg)
}
func (myLog *MyLog) Warning(eid uint32, msg string) error {
	if !IsDebug {
		return nil
	}
	return elog.Warning(eid, msg)
}
func (myLog *MyLog) Error(eid uint32, msg string) error {
	if !IsDebug {
		return nil
	}
	return elog.Error(eid, msg)
}

func (MyLog *MyLog) Close() {
	if elog != nil {
		elog.Close()
	}
}

var elog *eventlog.Log

func init() {
	elog, _ = eventlog.Open("test")
}
