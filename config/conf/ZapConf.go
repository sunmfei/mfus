package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type LogConfig struct {
	Level      string `json:"level"`                    // Level 最低日志等级，DEBUG<INFO<WARN<ERROR<FATAL 例如：info-->收集info等级以上的日志
	FileName   string `json:"file_name"`                // FileName 日志文件位置
	MaxSize    int    `json:"max_size"`                 // MaxSize 进行切割之前，日志文件的最大大小(MB为单位)，默认为100MB
	MaxAge     int    `json:"max_age"`                  // MaxAge 是根据文件名中编码的时间戳保留旧日志文件的最大天数。
	MaxBackups int    `json:"max_backups"`              // MaxBackups 是要保留的旧日志文件的最大数量。默认是保留所有旧的日志文件（尽管 MaxAge 可能仍会导致它们被删除。）
	Compress   bool   `json:"compress" yaml:"compress"` // Compress 确定是否应使用 gzip 压缩旋转的日志文件。默认是不执行压缩。
}

func NewLogConfig() *LogConfig {

	return &LogConfig{
		Level:      viper.GetString("logger.Level"),
		FileName:   fmt.Sprintf("%s-%v.log", viper.GetString("logger.FileName"), time.Now().Format("20060102")),
		MaxSize:    viper.GetInt("logger.MaxSize"),
		MaxAge:     viper.GetInt("logger.MaxAge"),
		MaxBackups: viper.GetInt("logger.MaxBackups"),
		Compress:   viper.GetBool("logger.Compress"),
	}
}
