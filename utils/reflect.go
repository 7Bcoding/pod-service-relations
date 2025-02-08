package utils

import "reflect"

// nested copy not supported
func CopyFields(src interface{}, dst interface{}) {
	tDst := reflect.TypeOf(dst).Elem()
	vDst := reflect.ValueOf(dst).Elem()
	tSrc := reflect.TypeOf(src)
	vSrc := reflect.ValueOf(src)
	if tSrc.Kind() == reflect.Ptr {
		tSrc = tSrc.Elem()
		vSrc = vSrc.Elem()
	}
	for i := 0; i < tSrc.NumField(); i++ {
		fName := tSrc.Field(i).Name
		if _, ok := tDst.FieldByName(fName); ok {
			v := vDst.FieldByName(fName)
			if v.Kind() == vSrc.Field(i).Kind() {
				v.Set(vSrc.Field(i))
			}
		}
	}
}
