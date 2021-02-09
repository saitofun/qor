package validations_test

import (
	"fmt"
	"testing"

	"github.com/saitofun/qor/validations"
)

func TestNewError(t *testing.T) {
	e := validations.NewError(&User{}, "Name", "Empty")
	fmt.Println(e.Error())
	fmt.Println(e.(*validations.Error).Label())
}
