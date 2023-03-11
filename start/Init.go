package start

import (
	"github.com/sunmfei/mfus/utils/commonUtils"
	"github.com/sunmfei/mfus/utils/dbutils"
	"github.com/sunmfei/mfus/utils/dbutils/redisUtil"
	"gorm.io/gorm"
)

func Page(workDir string, typ string, info interface{}) (*gorm.DB, *redisUtil.RedisPlay) {
	Setup(workDir)
	var db *gorm.DB
	db = GetDB(typ)
	autoTable(db, info)
	return db, GetRedis()
}

func Setup(workDir string) {
	//读取配置文件
	commonUtils.ReadConf(&workDir)
	initLogger()

}
func AutoTable(db *gorm.DB, info interface{}) {
	autoTable(db, info)
}
func GetDB(typ string) *gorm.DB {
	return dbutils.Setup(typ)
}

func GetRedis() *redisUtil.RedisPlay {
	return redisUtil.Setup()
}
