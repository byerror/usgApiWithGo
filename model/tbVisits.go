package model

import (
	"gorm.io/gorm"
)

// 这个没用到
type TbVisits struct {
	gorm.Model
	Ip string
}
