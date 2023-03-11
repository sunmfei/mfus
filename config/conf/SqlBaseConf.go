package conf

import "github.com/spf13/viper"

const (
	Oracle   = "oracle"
	Mysql    = "mysql"
	Postgres = "postgres"
)

type DatabaseConnConf struct {
	Type            string `json:"type" yaml:"type"`
	Url             string `json:"url" yaml:"url"`
	MaxOpenConn     string `json:"maxOpenConn" yaml:"maxOpenConn"`
	MaxIdleConn     string `json:"maxIdleConn" yaml:"maxIdleConn"`
	ConnMaxLifetime string `json:"connMaxLifetime" yaml:"connMaxLifetime"`
}

func NewDatabaseConnConf(typ string) *DatabaseConnConf {

	return &DatabaseConnConf{
		Type:            typ,
		Url:             viper.GetString("database." + typ + ".url"),
		MaxOpenConn:     viper.GetString("database." + typ + ".maxOpenConn"),
		MaxIdleConn:     viper.GetString("database." + typ + ".maxIdleConn"),
		ConnMaxLifetime: viper.GetString("database." + typ + ".connMaxLifetime"),
	}
}
