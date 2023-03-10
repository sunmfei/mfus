package dbutils

import (
	"github.com/sunmfei/mfus/common/MFei"
	"github.com/sunmfei/mfus/common/sun"
	"github.com/sunmfei/mfus/utils/dbutils/redisUtil"
)

// GetDB 开放给外部获得db连接
func GetDB(typ string) {
	MFei.DB = setup(typ)
	sun.Redis = redisUtil.Setup()
}
