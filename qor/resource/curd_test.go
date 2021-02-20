package resource_test

import (
	"fmt"
	"testing"

	"github.com/saitofun/qor/gorm"
	"github.com/saitofun/qor/qor"
	"github.com/saitofun/qor/qor/resource"
	"github.com/saitofun/qor/utils/test_db"
)

func TestResource_ToPrimaryQueryParams(t *testing.T) {
	type User struct {
		gorm.Model
		UserID string `gorm:"primaryKey"`
		Name   string
		Name2  *string
	}
	user := &User{}

	db := test_db.NewTestDB()
	db.AutoMigrate(&User{})

	ctx := &qor.Context{
		Request:     nil,
		Writer:      nil,
		CurrentUser: nil,
		Roles:       nil,
		ResourceID:  "",
		DB:          db,
		Config:      &qor.Config{DB: db},
		Errors:      qor.Errors{},
	}

	res := resource.New(user)
	sql, args := res.ToPrimaryQueryParams("1,1", ctx)
	fmt.Println(sql, args)
}
