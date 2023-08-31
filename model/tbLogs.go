package model

import (
	"gorm.io/gorm"
)

type TbLogs struct {
	gorm.Model
	Ip     string
	Action string //操作
	Result string //结果
	// CreateTime time.Time
}
