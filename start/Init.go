package start

import (
	"github.com/sunmfei/mfus/utils/commonUtils"
	"github.com/sunmfei/mfus/utils/dbutils"
	"github.com/sunmfei/mfus/utils/dbutils/redisUtil"
	"gorm.io/gorm"
)

func Page(workDir string, typ string, info interface{}) (*gorm.DB, *redisUtil.RedisPlay, error) {
	err := ReadConf(workDir)
	if err != nil {
		return nil, nil, err
	}
	var db *gorm.DB
	db, err = GetDB(typ)
	if err != nil {
		return nil, nil, err
	}
	redis, err := GetRedis()
	if err != nil {
		return nil, nil, err
	}
	err = AutoTable(db, info)
	if err != nil {
		return nil, nil, err
	}
	return db, redis, err
}
func ReadConf(workDir string) error {
	//读取配置文件
	return commonUtils.ReadConf(&workDir)
}
func AutoTable(db *gorm.DB, info interface{}) error {
	return autoTable(db, info)
}
func GetDB(typ string) (*gorm.DB, error) {
	return dbutils.Setup(typ)
}

func GetRedis() (*redisUtil.RedisPlay, error) {
	return redisUtil.Setup()
}
