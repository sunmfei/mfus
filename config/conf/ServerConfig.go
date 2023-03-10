package conf

import "github.com/spf13/viper"

type ServerConf struct {
	Host string `json:"host" yaml:"host" `
	Port string `json:"port" yaml:"port"`
}

func NewServerConf() *ServerConf {

	return &ServerConf{
		Host: viper.GetString("server.host"),
		Port: viper.GetString("server.port"),
	}
}
