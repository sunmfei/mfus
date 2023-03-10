package start

import (
	"github.com/sunmfei/mfus/utils"
	"github.com/sunmfei/mfus/utils/dbutils"
	"golang.org/x/exp/slog"
)

func Setup(workDir string, dbType string, info interface{}) {
	//读取配置文件
	utils.ReadConf(workDir)
	//日志文件配置
	err := utils.InitLogger()
	if err != nil {
		slog.Error("日志配置错误", err)
		return
	}
	//配置数据库
	dbutils.GetDB(dbType)

	//自动建表
	dbutils.AutoTable(info)

}
