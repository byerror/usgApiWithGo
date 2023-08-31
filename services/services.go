package services

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
	"yongyou/global"
	"yongyou/server"
	"yongyou/utils"

	"go.uber.org/zap"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/mgr"
)

var svcName = "usgApi"

// 安装服务
func InstallService() {

	// 打开服务管理器
	m, err := mgr.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer m.Disconnect()

	// s, _err := m.OpenService(svcName)
	// if _err != nil {
	// 	log.Printf("service %s already exists>>%s", svcName, _err.Error())
	// 	return
	// }

	// 获取当前可执行文件的路径
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
		return
	}

	//创建服务配置
	s, err := m.CreateService(svcName, exePath, mgr.Config{DisplayName: svcName, StartType: mgr.StartAutomatic})
	if err != nil {
		global.ZLog.Error("CreateService err", zap.Error(err))
		return
	}
	defer s.Close()
	global.ZLog.Info("安装服务成功，正在启动...")
	err = exec.Command("sc", "start", svcName).Start()
	if err != nil {
		global.ZLog.Info("服务启动失败!", zap.Error(err))
	} else {
		global.ZLog.Info("系统服务启动成功!")
	}

}

// 卸载服务
func Uninstall() {
	exec.Command("sc", "stop", svcName).Start()
	// 打开服务管理器
	m, err := mgr.Connect()
	if err != nil {
		global.ZLog.Error("连接服务管理器失败", zap.Error(err))
		return
	}
	defer m.Disconnect()

	// 打开指定名称的服务
	s, err := m.OpenService(svcName)
	if err != nil {
		global.ZLog.Error("打开指定服务失败", zap.Error(err))
		return
	}
	defer s.Close()

	// 删除服务
	err = s.Delete()
	if err != nil {
		global.ZLog.Error("删除服务失败", zap.Error(err))
	}
	global.ZLog.Info("服务移除成功！")
}

type exampleService struct{}

var elog *utils.MyLog

func (m *exampleService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
	// fasttick := time.Tick(500 * time.Millisecond)
	// slowtick := time.Tick(2 * time.Second)
	// tick := fasttick
	elog.Info(88, "start me")
	global.ZLog.Error("service .................start....................")
	go Start()
	elog.Error(444, "end me")
	global.ZLog.Error("service .................WS End....................")
loop:
	for {
		select {

		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				// Testing deadlock from https://code.google.com/p/winsvc/issues/detail?id=4
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				// golang.org/x/sys/windows/svc.TestExample is verifying this output.
				testOutput := strings.Join(args, "-")
				testOutput += fmt.Sprintf("-%d", c.Context)

				break loop
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}

			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

			default:

				global.ZLog.Error(fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func RunService(isDebug bool) {
	defer elog.Close()
	elog.Info(88, "run services...")
	name := svcName
	var err error
	global.ZLog.Info(fmt.Sprintf("starting %s service", name))
	run := svc.Run
	if isDebug {
		run = debug.Run
	}
	err = run(name, &exampleService{})
	if err != nil {
		elog.Error(44, "run service faile"+err.Error())
		global.ZLog.Error(fmt.Sprintf("%s service failed: %v", name, err))
		return
	}

	global.ZLog.Error(fmt.Sprintf("%s service stopped", name))
}

func Start() {
	defer func() {
		if err := recover(); err != nil {
			elog.Error(8884, fmt.Sprintf("start failed %v", err))
		}
	}()
	elog.Info(881, "pre LoadConfig ..")
	//1.初始始化配置
	if err := utils.LoadConfig(); err != nil {
		elog.Error(444, "load config err"+err.Error())

		global.ZLog.Error("load config err", zap.Error(err))
		return
	}
	elog.Info(882, "initDb ..")
	//2.加载数据库
	err := utils.InitDB()
	if err != nil {

		elog.Error(444, "err db"+err.Error())
		global.ZLog.Error("err db", zap.Error(err))
		return
	}
	elog.Info(882, "server.Start ..")
	global.ZLog.Info("init db ok!")
	//3.开启服务
	err = server.Start()
	if err != nil {
		elog.Error(444, "web server.start() err"+err.Error())

		global.ZLog.Error("err webServer", zap.Error(err))
	}
	elog.Info(888, "web server start OK")

}
