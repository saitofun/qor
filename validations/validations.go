package validations

import (
	"fmt"

	"github.com/saitofun/qor/gorm"
)

// NewError generate a new error for a model's field
// db    db context
// res   model value reference
// field invalid struct model filed name
// err   error message
func NewError(schema *gorm.Schema, res interface{}, field, err string) error {
	return &Error{
		schema:   schema,
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
	// scope := gorm.Scope{Value: err.Resource}
	// return fmt.Sprintf("%v_%v_%v", scope.GetModelStruct().ModelType.Name(), scope.PrimaryKeyValue(), err.Column)
	var (
		name = e.schema.ModelType.Name()
		id   = e.schema.LookUpField(name).ValueOf
	)
}

// Error show error message
func (e Error) Error() string {
	return fmt.Sprintf("%v", e.Message)
}
