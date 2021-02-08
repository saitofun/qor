package gorm_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/saitofun/qor/gorm"
)

type User struct {
	gorm.Model
	Name string
	Age  int `gorm:"check:age>0"`
}

func TestModelToSchema(t *testing.T) {
	var user User
	user.ID = 100
	schema, err := gorm.ModelToSchema(user)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("schema.Name:\n\t%s\n", schema.Name)
	fmt.Printf("schema.Table:\n\t%s\n", schema.Table)
	fmt.Printf("schema.ModelType:\n\t%v\n", schema.ModelType.String())
	fmt.Println("schema.DBNames")
	for _, v := range schema.DBNames {
		fmt.Println("\t", v)
	}
	fmt.Println("schema.PrimaryFieldDBNames")
	for i, v := range schema.PrimaryFields {
		fmt.Printf("\tindex:   %v\n", i)
		fmt.Printf("\tname:    %v\n", v.Name)
		fmt.Printf("\ttag:     %v\n", v.Tag)
		fmt.Printf("\tdb_name: %v\n", v.DBName)
		fmt.Printf("\tvalue:   %v\n", v.ReflectValueOf(reflect.ValueOf(user)))
	}
	fmt.Printf("schema.PrioritizedPrimaryField: \n\t%s\n",
		schema.PrioritizedPrimaryField.Name)
	fmt.Printf("schema.PrioritizedPrimaryField's Value: \n\t%s\n",
		schema.PrioritizedPrimaryField.ReflectValueOf(reflect.ValueOf(user)))
}
