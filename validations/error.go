package validations

import (
	"fmt"
	"reflect"

	"github.com/saitofun/qor/gorm"
)

// NewError generate a new error for a model's field
func NewError(resource interface{}, column, err string) error {
	return &Error{Resource: resource, Column: column, Message: err}
}

// Error is a validation error struct that hold model, column and error message
type Error struct {
	Resource interface{}
	Column   string
	Message  string
}

// Label is a label including model type, primary key and column name
func (e Error) Label() string {
	schema, _ := gorm.Parse(e.Resource)
	primary, _ := schema.PrioritizedPrimaryField.ValueOf(reflect.ValueOf(e.Resource))
	return fmt.Sprintf("%v_%v_%v", schema.ModelType.Name(), primary, e.Column)
}

// Error show error message
func (e Error) Error() string {
	return fmt.Sprintf("%v", e.Message)
}
