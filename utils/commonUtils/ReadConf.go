package commonUtils

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

// ReadConf 读取配置文件配置
func ReadConf(workDir *string) {
	if workDir == nil {
		return
	}
	//workDir, _ := os.Getwd()
	//log.Println("workDir：", workDir)
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(*workDir)
	err := viper.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			slog.Error("找不到配置文件。。。", err)
		} else {
			slog.Error("配置文件出错。。。", err)
		}
	}
	//打印获取到的配置文件key
	slog.Info("打印获取到的配置文件key :", viper.AllKeys())
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		slog.Info("Config file changed:", e.Name)
	})

}
