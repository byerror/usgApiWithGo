package utils

import (
	"fmt"
	"yongyou/global"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func LoadConfig() (err error) {
	v := viper.New()
	v.SetConfigFile(GetExeDir() + "\\config.yml")
	v.SetConfigType("yaml")
	err = v.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Fatal error config file: %s", err)

	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		global.ZLog.Info("config file changed:" + e.Name)
		if err = v.Unmarshal(&global.Cfg); err != nil {
			global.ZLog.Error("OnConfigChange err ", zap.Error(err))
		}
	})
	if err = v.Unmarshal(&global.Cfg); err != nil {
		return err
	}

	return nil

}
