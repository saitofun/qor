package gorm

import (
	"sync"

	"gorm.io/gorm/schema"
)

type (
	Field       = schema.Field
	StructField = schema.Field
)

func Parse(value interface{}) (*schema.Schema, error) {
	return schema.Parse(value, &sync.Map{}, schema.NamingStrategy{})
}
