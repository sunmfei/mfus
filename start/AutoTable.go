package start

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"reflect"
)

// AutoTable 自动建表
func autoTable(DB *gorm.DB, info interface{}) error {
	if info == nil || DB == nil {
		return nil
	}
	//1、扫描po包获取对象
	// 检测User结构体对应的表是否存在

	//MFei.LOGGER.Info("开始扫描需要创建的表...")
	typ := reflect.TypeOf(info)
	num := typ.NumField()
	if num <= 0 {
		return errors.New("未扫描岛所需要创建的表")
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
				fmt.Println("创建"+v+"失败...", err)
				continue
			}
			for i := 0; i < structTyp.NumField(); i++ {
				fmt.Println("", structTyp.Field(i).Name)
			}

			fmt.Println(v + "创建成功...")
		}

	}
	//MFei.LOGGER.Info("所有表均已创建...")
	return nil
}
