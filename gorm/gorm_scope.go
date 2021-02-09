package gorm

import (
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Scope struct {
	*gorm.Statement
	*schema.Schema
	val interface{}
}

func (s *Scope) New(v interface{}) *Scope {
	return s
}

func (s *Scope) IndirectValue() reflect.Value {
	return indirect(s.val)
}
