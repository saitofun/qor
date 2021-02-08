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
	DB        = gorm.DB
	Model     = gorm.Model // Model based columns: autoincrement id and time fields
	Config    = gorm.Config
	Statement = gorm.Statement

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
