package start

import (
	"github.com/sunmfei/mfus/common/MFei"
	"gorm.io/gorm"
	"reflect"
)

// AutoTable 自动建表
func autoTable(DB *gorm.DB, info interface{}) {
	if info == nil || DB == nil {
		return
	}
	//1、扫描po包获取对象
	// 检测User结构体对应的表是否存在

	//MFei.LOGGER.Info("开始扫描需要创建的表...")
	typ := reflect.TypeOf(info)
	num := typ.NumField()
	if num <= 0 {
		MFei.LOGGER.Info("未扫描岛所需要创建的表...")
		return
	}

	for i := 0; i < num; i++ {

		structTyp := typ.Field(i).Type
		//fmt.Println("类型为：", structTyp)
		class := reflect.New(structTyp)

		val := class.MethodByName("TableName").Call(make([]reflect.Value, 0))
		v := val[0].String()
		//MFei.LOGGER.Info("开始创建" + v + "...")
		if !DB.Migrator().HasTable(v) { //判断表是否存在
			//MFei.sun.Migrator().HasTable(reflect.New()) //自动创建
			err := DB.Migrator().CreateTable(class.Interface())
			if err != nil {
				MFei.LOGGER.Error("创建"+v+"失败...", err)
				continue
			}
			for i := 0; i < structTyp.NumField(); i++ {
				MFei.LOGGER.Info("", structTyp.Field(i).Name)
			}

			MFei.LOGGER.Info(v + "创建成功...")
		}

	}
	//MFei.LOGGER.Info("所有表均已创建...")

}
