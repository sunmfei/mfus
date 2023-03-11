package commonUtils

import "reflect"

// Copy 对象copy
func Copy(dis interface{}, src interface{}) {

	srcTyp := reflect.TypeOf(src).Elem()
	disTyp := reflect.TypeOf(dis).Elem()

	srcV := reflect.ValueOf(src).Elem()
	disV := reflect.ValueOf(dis).Elem()

	num := disTyp.NumField()
	for i := 0; i < num; i++ {
		name := disTyp.Field(i).Name
		srcSf, sb := srcTyp.FieldByName(name)
		disSf, db := disTyp.FieldByName(name)
		if !sb || !db {
			continue
		}

		if srcSf.Type == disSf.Type {
			disV.FieldByName(name).Set(srcV.FieldByName(name))
		}

	}

}
