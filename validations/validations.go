package validations

import (
	"fmt"

	"github.com/saitofun/qor/gorm"
)

// NewError generate a new error for a model's field
func NewError(res interface{}, field, err string) error {
	return &Error{
		Resource: res,
		Column:   field,
		Message:  err,
	}
}

// Error is a validation error struct that hold model, column and error message
type Error struct {
	schema   *gorm.Schema
	Resource interface{}
	Column   string
	Message  string
}

// Label is a label including model type, primary key and column name
func (e Error) Label() string {
	schema, _ := gorm.ModelToSchema(e.Resource)
	return fmt.Sprintf("%v_%v_%v",
		schema.ModelType.Name(),
		gorm.ReflectFieldValue(e.Resource, schema.PrioritizedPrimaryField),
		e.Column)
}

// Error show error message
func (e Error) Error() string {
	return fmt.Sprintf("%v", e.Message)
}
