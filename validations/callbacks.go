package validations

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/saitofun/qor/gorm"
)

var skipValidations = "validations:skip_validations"

func flatValidatorErrors(err govalidator.Errors) (ret []govalidator.Error) {
	for _, v := range err.Errors() {
		if errors, ok := v.(govalidator.Errors); ok {
			for _, e := range errors {
				ret = append(ret, e.(govalidator.Error))
			}
		} else if e, ok := v.(govalidator.Error); ok {
			ret = append(ret, e)
		}
	}
	return
}

func formattedError(err govalidator.Error, res interface{}) error {
	var (
		msg  = err.Error()
		name = err.Name
	)

	if strings.Index(msg, "non zero value required") >= 0 {
		msg = fmt.Sprintf("%v can't be blank", name)
	} else if strings.Index(msg, "as length") >= 0 {
		reg, _ := regexp.Compile(`\(([0-9]+)\|([0-9]+)\)`)
		sub := reg.FindSubmatch([]byte(err.Error()))
		msg = fmt.Sprintf("%v is the wrong length (should be %v~%v characters)",
			name, string(sub[1]), string(sub[2]))
	} else if strings.Index(msg, "as numeric") >= 0 {
		msg = fmt.Sprintf("%v is not a number", name)
	} else if strings.Index(msg, "as email") >= 0 {
		msg = fmt.Sprintf("%v is not a valid email address", name)
	}
	return NewError(res, name, msg)
}

// RegisterCallbacks register callback into GORM DB
func RegisterCallbacks(db *gorm.DB) {
	db.Callback().Create().Before("gorm:create").
		Register("validations:validate", validate)
	db.Callback().Update().Before("gorm:update").
		Register("validations:validate", validate)
}

type Validator interface {
	Validate(*gorm.DB)
}

func validate(db *gorm.DB) {
	if _, ok := db.Get("gorm:update_column"); ok {
		return
	}
	result, ok := db.Get(skipValidations)
	if ok && result.(bool) {
		return
	}
	if db.Error != nil {
		return
	}

	val := reflect.ValueOf(db.Statement.Model)
	model := val.Interface()
	if model == nil {
		_ = db.AddError(errors.New("model value is nil"))
		return
	}

	if v, ok := model.(Validator); ok {
		v.Validate(db)
		if db.Error != nil {
			return
		}
	}
	_, err := govalidator.ValidateStruct(model)
	db.Statement.ReflectValue = reflect.ValueOf(model).Elem()
	if err == nil {
		return
	}
	errs, ok := err.(govalidator.Errors)
	if ok {
		for _, e := range flatValidatorErrors(errs) {
			_ = db.AddError(formattedError(e, model))
		}
	} else {
		_ = db.AddError(err)
	}
}
