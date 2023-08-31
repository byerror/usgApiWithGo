package server

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"
	"yongyou/global"
	"yongyou/model"
	"yongyou/utils"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code   int    // 正常返回0 ，当获取不到策略状态时，返回-1
	Status bool   //当前策略状态 开启时 表示已断开，关闭时表示可连接
	Now    string //时间
	Msg    string //操作消息 查询成功，开启/关闭操作成功

	DayPV         int
	TotalPV       int
	OperTimes     int
	TbLogs        []model.TbLogs
	ServiceStatus bool //服务状态
}

// 看是否有cookie,没的话就设置,返回首次访问true/false
func Get2SetCookie(ctx *gin.Context) (firstVist bool) {
	_, err := ctx.Cookie("pv")

	if err == nil { //存在，表示当天访问过。
		return false
	}

	//当天首次访问
	ctx.SetCookie("pv", "1", 1*24*3600, "/", "", false, true)
	return true

}

func Start() (err error) {

	//3.开启服务
	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"FormatAsDate": utils.FormatAsDate,
	})
	router.LoadHTMLGlob(utils.GetExeDir() + "/views/*")
	//后面的是物理路径，替换成左边的虚拟路径
	router.Static("/staticx", utils.GetExeDir()+"/static")

	router.NoRoute(func(ctx *gin.Context) {
		// ctx.Redirect(http.StatusNotFound, "/404.html")

		ctx.JSON(http.StatusNotFound, gin.H{"message": "Not Found " + ctx.Request.URL.Path})
	})
	router.GET("/", _GET)

	router.POST("/", _POST)

	err = router.Run(fmt.Sprintf(":%d", global.Cfg.Server.Port)) //127.0.0.1:9999")
	if err != nil {
		return err
	} else {
		return nil
	}

}

// 查询
func _GET(ctx *gin.Context) {
	ret := &Response{Msg: "查询成功", Code: 0, Now: time.Now().Format("15:04:05"), Status: true}
	// buff, _ := json.Marshal(ret)
	if Get2SetCookie(ctx) {
		utils.InserLogs(ctx.ClientIP(), "访问", "-")
	}
	//获取日志查询消息
	ret.TbLogs = utils.GetLogs()

	//获取报表
	report := utils.GetReport()
	ret.DayPV = report.DayPV
	ret.TotalPV = report.TotalPV
	ret.OperTimes = report.OperTimes
	//获取策略状态
	isEnable, err := utils.GetPolicy()
	if err != nil {
		ret.Code = -1
		ret.Status = false
		ret.Msg = "获取策略状态失败，请联系IT人员>>" + err.Error()
	} else {
		ret.Status = isEnable
	}
	//获取服务状态
	isOpen, _ := utils.CheckPorts(global.Cfg.Monitor.IP, []int{global.Cfg.Monitor.Port})
	ret.ServiceStatus = isOpen

	ctx.HTML(http.StatusOK, "index.html", ret)

}

// 操作
func _POST(ctx *gin.Context) {
	action, _ := ctx.GetPostForm("action") //前端发来的当前状态
	//ret := &Response{Msg: "查询成功", Now: time.Now().Format("2006-01-02 15:04:05"), Status: true}

	//原来的状态  tue为已阻断
	isEnable := strings.EqualFold(action, "true")
	//换成新的状态则要反过来
	sAction := "阻断"
	if isEnable {
		sAction = "开通"
	}
	err := utils.SetPolicy(!isEnable)

	if err != nil {
		utils.InserLogs(ctx.ClientIP(), sAction, err.Error())
		ctx.JSON(http.StatusOK, gin.H{"code": -1})
	} else {
		utils.InserLogs(ctx.ClientIP(), sAction, "成功")
		ctx.JSON(http.StatusOK, gin.H{"code": 0})
	}

}
