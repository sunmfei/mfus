package conf

import (
	"github.com/spf13/viper"
	"time"
)

type RedisConf struct {
	Host        string        `json:"host" yaml:"host"`
	Password    string        `json:"password" yaml:"password"`
	DB          int           `json:"sun" yaml:"sun"`
	IdleTimeout time.Duration `json:"idleTimeout" yaml:"idleTimeout"`
}

func NewRedisConf() *RedisConf {

	return &RedisConf{
		Host:     viper.GetString("database.redis.host"),
		Password: viper.GetString("database.redis.password"),
		DB:       viper.GetInt("database.redis.sun"),
	}
}
