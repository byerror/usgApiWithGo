package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"yongyou/global"
	"yongyou/services"
	"yongyou/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sys/windows/svc"
)

var elog *utils.MyLog

func main() {

	defer elog.Close()

	iniZap()

	// global.ZLog, _ = zap.NewDevelopment()
	//服务判断
	isService, err := svc.IsWindowsService()
	if err != nil {
		global.ZLog.Error("failed to determine if we are running in service:", zap.Error(err))
		return
	}
	if isService {
		global.ZLog.Info("启动服务中...")
		services.RunService(false)
		return
	}

	if len(os.Args) < 2 {
		global.ZLog.Error("指定参数,可用：install,uninstall")
		services.InstallService()
		return
	}
	cmd := strings.ToLower(os.Args[1])
	if cmd == "install" {
		services.InstallService()
	} else if cmd == "uninstall" {
		services.Uninstall()
	}

}

func iniZap() {
	elog.Info(8, "iniZap...")
	defer func() {
		if err := recover(); err != nil {
			elog.Info(4, fmt.Sprintf("err iniZap %v", err))
		}
	}()
	isDebug := true
	if global.Cfg == nil || global.Cfg.Server.Debug == nil {
		//global.Cfg.Server.Debug = &isDebug
	} else {
		isDebug = *global.Cfg.Server.Debug
	}
	elog.Info(8, "zap.Config pre...")
	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:       isDebug,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			TimeKey:        "time",
			NameKey:        "logger",
			CallerKey:      "file",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
		OutputPaths:      []string{"stdout", utils.GetExeDir() + "/zap.log"},
		ErrorOutputPaths: []string{"stderr", utils.GetExeDir() + "/zap_err.log"},
		InitialFields: map[string]interface{}{
			"app": "test",
		},
	}

	var err error
	global.ZLog, err = cfg.Build()
	if err != nil {
		log.Fatalln("zlog err" + err.Error())
	}

}
