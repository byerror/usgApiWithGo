package utils

import (
	"fmt"
	"yongyou/global"
	"yongyou/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Reports struct {
	DayPV     int
	TotalPV   int
	OperTimes int
}

func InitDB() error {
	dial := sqlite.Open(global.Cfg.Database.Path)
	var err error
	global.Db, err = gorm.Open(dial, &gorm.Config{})
	if err != nil {
		return err
	}
	global.Db.AutoMigrate(&model.TbLogs{})
	global.Db.AutoMigrate(&model.TbVisits{})
	return nil

}

// func Test() {
// 	dial := sqlite.Open("db.db")
// 	db, err := gorm.Open(dial, &gorm.Config{})
// 	if err != nil {
// 		fmt.Println("err>>", err.Error())
// 		return
// 	}
// 	// Migrate the schema
// 	db.AutoMigrate(&model.TbLogs{})
// 	db.AutoMigrate(&model.TbVisits{})

// 	db.Create(&model.TbLogs{Ip: "222"})
// 	var tbLogs []model.TbLogs
// 	db.Where("id > ?", 0).Find(&tbLogs)
// 	fmt.Print(tbLogs)

// }

// 保存日志 action:访问，关闭，开启
func InserLogs(ip string, action string, result string) {
	var tbLogs = model.TbLogs{Ip: ip, Action: action, Result: result}
	global.Db.Create(&tbLogs)
}

// 获取日志
func GetLogs() (tbLogs []model.TbLogs) {
	global.Db.Order("id DESC").Limit(10).Find(&tbLogs)
	return tbLogs
}

// 获取访问及操作报表
func GetReport() (report Reports) {

	var reports Reports
	sql := `select sum(case when  strftime('%Y-%m-%d',created_at)=date() and action='访问' then 1 else 0 end) as dayPV,
		sum(case when action='访问' then 1 else 0 end ) as totalPV,
		sum(case when action='访问' then 0 else 1 end) as operTimes 
		from tb_logs 
		`
	row := global.Db.Raw(sql).Row() //.Scan(&reports) //, &reports)
	row.Scan(&reports.DayPV, &reports.TotalPV, &reports.OperTimes)
	fmt.Println(row, reports)
	return reports
}
