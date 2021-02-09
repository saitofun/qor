package validations

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/saitofun/qor/gorm"
)

var (
	skipValidations = "validations:skip_validations"
	validProcKey    = "validations:validate"
)

func flatValidatorErrors(validatorErrors govalidator.Errors) []govalidator.Error {
	resultErrors := make([]govalidator.Error, 0)
	for _, validatorError := range validatorErrors.Errors() {
		if errors, ok := validatorError.(govalidator.Errors); ok {
			for _, e := range errors {
				resultErrors = append(resultErrors, e.(govalidator.Error))
			}
		}
		if e, ok := validatorError.(govalidator.Error); ok {
			resultErrors = append(resultErrors, e)
		}
	}
	return resultErrors
}

func formattedError(err govalidator.Error, res interface{}) error {
	msg := err.Error()
	attr := err.Name
	if strings.Index(msg, "non zero value required") >= 0 {
		msg = fmt.Sprintf("%v can't be blank", attr)
	} else if strings.Index(msg, "as length") >= 0 {
		reg, _ := regexp.Compile(`\(([0-9]+)\|([0-9]+)\)`)
		submatch := reg.FindSubmatch([]byte(err.Error()))
		msg = fmt.Sprintf("%v is the wrong length (should be %v~%v characters)",
			attr, string(submatch[1]), string(submatch[2]))
	} else if strings.Index(msg, "as numeric") >= 0 {
		msg = fmt.Sprintf("%v is not a number", attr)
	} else if strings.Index(msg, "as email") >= 0 {
		msg = fmt.Sprintf("%v is not a valid email address", attr)
	}
	return NewError(res, attr, msg)

}

// RegisterCallbacks register callbacks into GORM DB
func RegisterCallbacks(db *gorm.DB) {
	createProcessor := db.Callback().Create()
	updateProcessor := db.Callback().Update()

	if createProcessor.Get(validProcKey) == nil {
		createProcessor.Before("gorm:before_create").
			Register(validProcKey, validate)
	}
	if updateProcessor.Get(validProcKey) == nil {
		updateProcessor.Before("gorm:before_update").
			Register(validProcKey, validate)
	}
}

type Validator interface {
	Validate(*gorm.DB)
}

func validate(db *gorm.DB) {
	tx := db.Session(&gorm.Session{})
	if _, ok := tx.Get("gorm:update_column"); ok {
		return
	}
	if v, ok := tx.Get(skipValidations); ok && v.(bool) {
		return
	}
	if tx.Error != nil {
		return
	}
	v, ok := tx.Statement.Model.(Validator)
	if !ok {
		return
	}
	v.Validate(tx)
	if tx.Error != nil {
		return
	}
	_, validatorErrors := govalidator.ValidateStruct(v)
	if validatorErrors != nil {
		if errors, ok := validatorErrors.(govalidator.Errors); ok {
			for _, err := range flatValidatorErrors(errors) {
				db.AddError(formattedError(err, v))
			}
		} else {
			db.AddError(validatorErrors)
		}
	}
}
