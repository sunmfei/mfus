package config

import "github.com/sunmfei/mfus/config/conf"

type Config struct {
	ServerConf *conf.ServerConf `json:"serverConf" yaml:"serverConf"`
	RedisConf  *conf.RedisConf  `json:"redisConf" yaml:"redisConf"`
}

func NewConfig() *Config {

	return &Config{
		ServerConf: conf.NewServerConf(),
		RedisConf:  conf.NewRedisConf(),
	}

}
