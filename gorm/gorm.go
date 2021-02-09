// author: birdyfj@gmail.com
// this gorm is wrapped all gorm2's defines and try to make QOR frame migrated to gorm2

package gorm

import (
	"reflect"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// gorm v2 type defines here
type (
	Model     = gorm.Model // Model based columns: autoincrement id and time fields
	Config    = gorm.Config
	Statement = gorm.Statement
	Session   = gorm.Session

	Field  = schema.Field
	Schema = schema.Schema
)

// gorm v2 func defines here
var (
	Open = gorm.Open
)

// gorm.logger.LogLevel
const (
	DBLogSilent logger.LogLevel = iota + 1
	DBLogError
	DBLogWarn
	DBLogInfo
)

// gorm error defines here
var (
	ErrRecordNotFound        = gorm.ErrRecordNotFound
	ErrInvalidTransaction    = gorm.ErrInvalidTransaction
	ErrNotImplemented        = gorm.ErrNotImplemented
	ErrMissingWhereClause    = gorm.ErrMissingWhereClause
	ErrUnsupportedRelation   = gorm.ErrUnsupportedRelation
	ErrPrimaryKeyRequired    = gorm.ErrPrimaryKeyRequired
	ErrModelValueRequired    = gorm.ErrModelValueRequired
	ErrInvalidData           = gorm.ErrInvalidData
	ErrUnsupportedDriver     = gorm.ErrUnsupportedDriver
	ErrRegistered            = gorm.ErrRegistered
	ErrInvalidField          = gorm.ErrInvalidField
	ErrEmptySlice            = gorm.ErrEmptySlice
	ErrDryRunModeUnsupported = gorm.ErrDryRunModeUnsupported
)

func ModelToSchema(model interface{}, db ...*gorm.DB) (*Schema, error) {
	var namer schema.Namer = schema.NamingStrategy{}
	if len(db) > 0 && db[0] != nil && db[0].Config != nil {
		namer = db[0].Config.NamingStrategy
	}
	return schema.Parse(model, &sync.Map{}, namer)
}

func ReflectFieldValue(model interface{}, field *Field) interface{} {
	return field.ReflectValueOf(reflect.Indirect(reflect.ValueOf(model))).Interface()
}

func ReflectIndirectFieldValue(model interface{}, field *Field) interface{} {
	return reflect.Indirect(reflect.ValueOf(ReflectFieldValue(model, field))).Interface()
}

func PrimaryFields(model interface{}) []*Field {
	if schema, err := ModelToSchema(model); err != nil {
		return nil
	} else {
		return schema.PrimaryFields
	}
}

func PrimaryField(model interface{}) *Field {
	if schema, err := ModelToSchema(model); err != nil {
		return nil
	} else {
		return schema.PrioritizedPrimaryField
	}
}

func PrimaryKeyZero(model interface{}) bool {
	schema, err := ModelToSchema(model)
	if err != nil {
		return false
	}
	return schema.PrioritizedPrimaryField == nil &&
		isBlank(reflect.ValueOf(ReflectFieldValue(model, schema.PrioritizedPrimaryField)))
}

func PrimaryKeyValue(model interface{}) interface{} {
	schema, _ := ModelToSchema(model)
	return ReflectFieldValue(model, schema.PrioritizedPrimaryField)
}

func IsFieldBlank(model interface{}, field *Field) bool {
	return isBlank(reflect.ValueOf(ReflectFieldValue(model, field)))
}

func isBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}
