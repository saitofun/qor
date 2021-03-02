package gorm

import (
	"reflect"
	"sync"
	"time"

	"gorm.io/gorm/schema"
)

type (
	Field       = schema.Field
	StructField = schema.Field
)

func Parse(value interface{}) (*schema.Schema, error) {
	return schema.Parse(value, &sync.Map{}, schema.NamingStrategy{})
}

func IsNormalField(f *Field) bool {
	typ := f.FieldType
	for typ.Kind() == reflect.Ptr || typ.Kind() == reflect.Slice || typ.Kind() == reflect.Interface {
		typ = typ.Elem()
	}
	switch typ.Kind() {
	case reflect.String,
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64:
		return true
	default:
		if typ.Kind() == reflect.Struct {
			if _, ok := reflect.Zero(typ).Interface().(time.Time); ok {
				return true
			}
		}
		return false
	}
}
