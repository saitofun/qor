package resource

import (
	"reflect"

	"github.com/saitofun/qor/gorm"
	"github.com/saitofun/qor/qor"
	"github.com/saitofun/qor/qor/utils"
	"github.com/saitofun/qor/roles"
)

// Resourcer interface
type Resourcer interface {
	GetResource() *Resource
	GetMetas([]string) []Metaor
	CallFindMany(interface{}, *qor.Context) error
	CallFindOne(interface{}, *MetaValues, *qor.Context) error
	CallSave(interface{}, *qor.Context) error
	CallDelete(interface{}, *qor.Context) error
	NewSlice() interface{}
	NewStruct() interface{}
}

// ConfigureResourceBeforeInitializeInterface if a struct implemented this interface, it will be called before everything when create a resource with the struct
type ConfigureResourceBeforeInitializeInterface interface {
	ConfigureQorResourceBeforeInitialize(Resourcer)
}

// ConfigureResourceInterface if a struct implemented this interface, it will be called after configured by user
type ConfigureResourceInterface interface {
	ConfigureQorResource(Resourcer)
}

// Resource is a struct that including basic definition of qor resource
type Resource struct {
	Name            string
	Value           interface{}
	PrimaryFields   []*gorm.Field
	FindManyHandler func(interface{}, *qor.Context) error
	FindOneHandler  func(interface{}, *MetaValues, *qor.Context) error
	SaveHandler     func(interface{}, *qor.Context) error
	DeleteHandler   func(interface{}, *qor.Context) error
	Permission      *roles.Permission
	Validators      []*Validator
	Processors      []*Processor
	primaryField    *gorm.Field
}

// New initialize qor resource
func New(value interface{}) *Resource {
	schema, _ := gorm.ModelToSchema(value)
	return &Resource{
		Value:           value,
		Name:            utils.HumanizeString(utils.ModelType(value).Name()),
		FindOneHandler:  (*Resource)(nil).findOneHandler,
		FindManyHandler: (*Resource)(nil).findManyHandler,
		SaveHandler:     (*Resource)(nil).saveHandler,
		DeleteHandler:   (*Resource)(nil).deleteHandler,
		PrimaryFields:   schema.PrimaryFields,
		primaryField:    schema.PrioritizedPrimaryField,
	}
}

// GetResource return itself to match interface `Resourcer`
func (res *Resource) GetResource() *Resource {
	return res
}

// Validator validator struct
type Validator struct {
	Name    string
	Handler func(interface{}, *MetaValues, *qor.Context) error
}

// AddValidator add validator to resource, it will invoked when creating, updating, and will rollback the change if validator return any error
func (res *Resource) AddValidator(validator *Validator) {
	for idx, v := range res.Validators {
		if v.Name == validator.Name {
			res.Validators[idx] = validator
			return
		}
	}

	res.Validators = append(res.Validators, validator)
}

// Processor processor struct
type Processor struct {
	Name    string
	Handler func(interface{}, *MetaValues, *qor.Context) error
}

// AddProcessor add processor to resource, it is used to process data before creating, updating, will rollback the change if it return any error
func (res *Resource) AddProcessor(processor *Processor) {
	for idx, p := range res.Processors {
		if p.Name == processor.Name {
			res.Processors[idx] = processor
			return
		}
	}
	res.Processors = append(res.Processors, processor)
}

// NewStruct initialize a struct for the Resource
func (res *Resource) NewStruct() interface{} {
	if res.Value == nil {
		return nil
	}
	return reflect.New(utils.Indirect(reflect.ValueOf(res.Value)).Type()).Interface()
}

// NewSlice initialize a slice of struct for the Resource
func (res *Resource) NewSlice() interface{} {
	if res.Value == nil {
		return nil
	}
	sliceType := reflect.SliceOf(reflect.TypeOf(res.Value))
	slice := reflect.MakeSlice(sliceType, 0, 0)
	slicePtr := reflect.New(sliceType)
	slicePtr.Elem().Set(slice)
	return slicePtr.Interface()
}

// GetMetas get defined metas, to match interface `Resourcer`
func (res *Resource) GetMetas([]string) []Metaor {
	panic("not defined")
}

// HasPermission check permission of resource
func (res *Resource) HasPermission(mode roles.PermissionMode, context *qor.Context) bool {
	if res == nil || res.Permission == nil {
		return true
	}

	var roles = make([]interface{}, 0)
	for _, role := range context.Roles {
		roles = append(roles, role)
	}
	return res.Permission.HasPermission(mode, roles...)
}
