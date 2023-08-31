package global

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ZLog *zap.Logger
	Db   *gorm.DB
	Cfg  *Config
)

type Config struct {
	Server Server `yaml:"server" json:"server"`
	Api    Api

	Database Database `yaml:"database" json:"database"`
	Monitor  Monitor
}

type Api struct {
	Ip     string
	User   string
	Psw    string
	Policy Policy `yaml:"policy" json:"policy"`
}
type Server struct {
	Debug *bool
	Port  int
}
type Database struct {
	Path string
}

type Policy struct {
	Rule     string
	Template string
}
type Monitor struct {
	IP   string
	Port int
}
