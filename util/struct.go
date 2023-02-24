package util

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
)

// CopyStruct 复制属性，src和dest必须是结构体指针
//func CopyStruct(dst, src any) (err error) {
//	// 防止意外panic
//	defer func() {
//		if e := recover(); e != nil {
//			err = errors.New(fmt.Sprintf("%v", e))
//		}
//	}()
//
//	dstType, dstValue := reflect.TypeOf(dst), reflect.ValueOf(dst)
//	srcType, srcValue := reflect.TypeOf(src), reflect.ValueOf(src)
//
//	// dst必须结构体指针类型
//	if dstType.Kind() != reflect.Ptr || dstType.Elem().Kind() != reflect.Struct {
//		return errors.New("dst type should be a struct pointer")
//	}
//
//	// src必须为结构体或者结构体指针
//	if srcType.Kind() == reflect.Ptr {
//		srcType, srcValue = srcType.Elem(), srcValue.Elem()
//	}
//	if srcType.Kind() != reflect.Struct {
//		return errors.New("src type should be a struct or a struct pointer")
//	}
//
//	// 取具体内容
//	dstType, dstValue = dstType.Elem(), dstValue.Elem()
//
//	// 属性个数
//	propertyNums := dstType.NumField()
//
//	for i := 0; i < propertyNums; i++ {
//		// 属性
//		property := dstType.Field(i)
//		// 待填充属性值
//		propertyValue := srcValue.FieldByName(property.Name)
//
//		// 无效，说明src没有这个属性 || 属性同名但类型不同
//		if !propertyValue.IsValid() || property.Type != propertyValue.Type() {
//			continue
//		}
//
//		if dstValue.Field(i).CanSet() {
//			dstValue.Field(i).Set(propertyValue)
//		}
//	}
//	return nil
//}

// StructToMap This function will help you to convert your object from struct to map[string]interface{} based on your JSON tag in your structs.
func StructToMap(src interface{}) map[string]interface{} {

	res := map[string]interface{}{}
	if src == nil {
		return res
	}
	v := reflect.TypeOf(src)
	reflectValue := reflect.ValueOf(src)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			//if v.Field(i).Type.Kind() == reflect.Struct {
			//res[tag] = structToMap(field)
			//	continue
			//} else {
			res[tag] = field
			//}
		}
	}
	return res
}

// CopyStruct 复制属性，src和dest必须是结构体指针
func CopyStruct[T any](src any) (dst *T, err error) {
	dst = new(T)
	// 防止意外panic
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
		}
	}()

	dstType, dstValue := reflect.TypeOf(dst), reflect.ValueOf(dst)
	srcType, srcValue := reflect.TypeOf(src), reflect.ValueOf(src)

	// dst必须结构体指针类型
	if dstType.Kind() != reflect.Ptr || dstType.Elem().Kind() != reflect.Struct {
		return nil, errors.New("dst type should be a struct pointer")
	}

	// src必须为结构体或者结构体指针
	if srcType.Kind() == reflect.Ptr {
		srcType, srcValue = srcType.Elem(), srcValue.Elem()
	}
	if srcType.Kind() != reflect.Struct {
		return nil, errors.New("src type should be a struct or a struct pointer")
	}

	// 取具体内容
	dstType, dstValue = dstType.Elem(), dstValue.Elem()

	// 属性个数
	propertyNums := dstType.NumField()

	for i := 0; i < propertyNums; i++ {
		// 属性
		property := dstType.Field(i)
		// 待填充属性值
		propertyValue := srcValue.FieldByName(property.Name)

		// 无效，说明src没有这个属性 || 属性同名但类型不同
		if !propertyValue.IsValid() || property.Type != propertyValue.Type() {
			continue
		}

		if dstValue.Field(i).CanSet() {
			dstValue.Field(i).Set(propertyValue)
		}
	}

	return dst, nil
}
