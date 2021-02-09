package gorm

import (
	"reflect"
)

func indirect(v interface{}) reflect.Value {
	val := reflect.ValueOf(v)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val
}
