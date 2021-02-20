// author: birdyfj@gmail.com
// this gorm is wrapped all gorm2's defines and try to make QOR frame migrated to gorm2

package gorm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// gorm v2 type defines here
type (
	Model        = gorm.Model // Model based columns: autoincrement id and time fields
	Statement    = gorm.Statement
	Session      = gorm.Session
	Schema       = schema.Schema
	Relationship = schema.Relationship
)

// gorm.logger.LogLevel
const (
	LogSilent logger.LogLevel = iota + 1
	LogError
	LogWarn
	LogInfo
)

const (
	HasOne    = schema.HasOne
	HasMany   = schema.HasMany
	BelongsTo = schema.BelongsTo
	Many2Many = schema.Many2Many
)

