package config

import "github.com/sunmfei/mfus/config/conf"

type Config struct {
	RedisConf *conf.RedisConf `json:"redisConf" yaml:"redisConf"`
}

func NewConfig() *Config {

	return &Config{
		RedisConf: conf.NewRedisConf(),
	}

}
