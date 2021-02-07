package gorm

import (
	gorm1 "github.com/jinzhu/gorm"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// gorm v2 type defines here
type (
	DB        = gorm.DB
	Model     = gorm.Model // Model based columns: autoincrement id and time fields
	Config    = gorm.Config
	Schema    = schema.Schema
	Statement = gorm.Statement
)

// gorm v2 func defines here
var (
	Open = gorm.Open
)

// gorm v2 enum defines here

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

type (
	_Model       = gorm1.Model
	_Scope       = gorm1.Scope
	_Field       = gorm1.Field
	_StructField = gorm1.StructField
)
