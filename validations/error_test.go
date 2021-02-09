package validations_test

import (
	"testing"

	"github.com/saitofun/qor/validations"
)

func TestNewError(t *testing.T) {
	e := validations.NewError(&User{}, "Name", "Empty")
	t.Logf(e.Error())
	label := e.(*validations.Error).Label()
	t.Logf(label)
}
