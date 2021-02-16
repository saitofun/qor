package resource_test

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"testing"

	"github.com/saitofun/qor/gorm"
	"github.com/saitofun/qor/qor"
	"github.com/saitofun/qor/qor/resource"
	"github.com/saitofun/qor/qor/utils"
	"github.com/saitofun/qor/utils/test_db"
)

func format(value interface{}) string {
	return fmt.Sprint(utils.Indirect(reflect.ValueOf(value)).Interface())
}

func checkMeta(record interface{}, meta *resource.Meta, value interface{}, t *testing.T, expectedValues ...string) {
	var (
		context       = &qor.Context{DB: test_db.TestDB()}
		metaValue     = &resource.MetaValue{Name: meta.Name, Value: value}
		expectedValue = fmt.Sprint(value)
	)

	for _, v := range expectedValues {
		expectedValue = v
	}

	meta.PreInitialize()
	meta.Initialize()

	if meta.Setter != nil {
		meta.Setter(record, metaValue, context)
		if context.HasError() {
			t.Errorf("No error should happen, but got %v", context.Errors)
		}

		result := meta.Valuer(record, context)
		if resultValuer, ok := result.(driver.Valuer); ok {
			if v, err := resultValuer.Value(); err == nil {
				result = v
			}
		}

		if format(result) != expectedValue {
			t.Errorf("Wrong value, should be %v, but got %v", expectedValue, format(result))
		}
	} else {
		t.Errorf("No setter generated for meta %v", meta.Name)
	}
}

func TestStringMetaValuerAndSetter(t *testing.T) {
	user := &struct {
		Name  string
		Name2 *string
	}{}

	res := resource.New(user)

	meta := &resource.Meta{
		Name:         "Name",
		BaseResource: res,
	}

	checkMeta(user, meta, "hello world", t)

	meta2 := &resource.Meta{
		Name:         "Name2",
		BaseResource: res,
	}

	checkMeta(user, meta2, "hello world2", t)
}

func TestIntMetaValuerAndSetter(t *testing.T) {
	user := &struct {
		Age  int
		Age2 uint
		Age3 *int8
		Age4 *uint8
	}{}

	res := resource.New(user)

	meta := &resource.Meta{
		Name:         "Age",
		BaseResource: res,
	}

	checkMeta(user, meta, 18, t)

	meta2 := &resource.Meta{
		Name:         "Age2",
		BaseResource: res,
	}

	checkMeta(user, meta2, "28", t)

	meta3 := &resource.Meta{
		Name:         "Age3",
		BaseResource: res,
	}

	checkMeta(user, meta3, 38, t)

	meta4 := &resource.Meta{
		Name:         "Age4",
		BaseResource: res,
	}

	checkMeta(user, meta4, "48", t)
}

func TestFloatMetaValuerAndSetter(t *testing.T) {
	user := &struct {
		Age  float64
		Age2 *float64
	}{}

	res := resource.New(user)

	meta := &resource.Meta{
		Name:         "Age",
		BaseResource: res,
	}

	checkMeta(user, meta, 18.5, t)

	meta2 := &resource.Meta{
		Name:         "Age2",
		BaseResource: res,
	}

	checkMeta(user, meta2, "28.5", t)
}

func TestBoolMetaValuerAndSetter(t *testing.T) {
	user := &struct {
		Actived  bool
		Actived2 *bool
	}{}

	res := resource.New(user)

	meta := &resource.Meta{
		Name:         "Actived",
		BaseResource: res,
	}

	checkMeta(user, meta, "true", t)

	meta2 := &resource.Meta{
		Name:         "Actived2",
		BaseResource: res,
	}

	checkMeta(user, meta2, "true", t)

	meta3 := &resource.Meta{
		Name:         "Actived",
		BaseResource: res,
	}

	checkMeta(user, meta3, "", t, "false")

	meta4 := &resource.Meta{
		Name:         "Actived2",
		BaseResource: res,
	}

	checkMeta(user, meta4, "f", t, "false")
}

type scanner struct {
	Body string
}

func (s *scanner) Scan(value interface{}) error {
	s.Body = fmt.Sprint(value)
	return nil
}

func (s scanner) Value() (driver.Value, error) {
	return s.Body, nil
}

func TestScannerMetaValuerAndSetter(t *testing.T) {
	user := &struct {
		Scanner scanner
	}{}

	res := resource.New(user)

	meta := &resource.Meta{
		Name:         "Scanner",
		BaseResource: res,
	}

	checkMeta(user, meta, "scanner", t)
}

func TestSliceMetaValuerAndSetter(t *testing.T) {
	t.Skip()

	user := &struct {
		Names  []string
		Names2 []*string
		Names3 *[]string
		Names4 []*string
	}{}

	res := resource.New(user)

	meta := &resource.Meta{
		Name:         "Names",
		BaseResource: res,
	}

	checkMeta(user, meta, []string{"name1", "name2"}, t)

	meta2 := &resource.Meta{
		Name:         "Names2",
		BaseResource: res,
	}

	checkMeta(user, meta2, []string{"name1", "name2"}, t)

	meta3 := &resource.Meta{
		Name:         "Names3",
		BaseResource: res,
	}

	checkMeta(user, meta3, []string{"name1", "name2"}, t)

	meta4 := &resource.Meta{
		Name:         "Names4",
		BaseResource: res,
	}

	checkMeta(user, meta4, []string{"name1", "name2"}, t)
}

type User struct {
	gorm.Model
	Name      string
	CompanyID int64
	Company   Company `gorm:"foreignKey:CompanyID"`
}

type Company struct {
	gorm.Model
	Name       string
	OrganizeID int64
	Organize   *Organize `gorm:"foreignKey:OrganizeID"`
}

type Organize struct {
	gorm.Model
	Name string
}

func TestUtil(t *testing.T) {
	var i, j int
	fmt.Println(reflect.ValueOf(i).Type() == reflect.ValueOf(j).Type())
}

func TestSetupValuer(t *testing.T) {
	db := test_db.TestDB()

	if err := db.AutoMigrate(&User{}, &Company{}, &Organize{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	user := &User{}
	meta := &resource.Meta{
		Name:         "Company.OrganizeID",
		BaseResource: resource.New(user),
	}
	ctx := &qor.Context{
		DB:     db,
		Config: &qor.Config{DB: db},
	}

	meta.PreInitialize()
	meta.Initialize()

	fmt.Println(meta.Valuer(&User{}, ctx))
	meta.Setter(user, &resource.MetaValue{
		Name:  "Company",
		Value: nil,
		Index: 0,
	}, ctx)
	fmt.Println(meta.Valuer(user, ctx))

	fmt.Println(meta.Valuer(&User{}, ctx))
	meta.Setter(user, &resource.MetaValue{
		Name:  "Company",
		Value: &Company{Name: "XXXX"},
		Index: 0,
	}, ctx)
	fmt.Println(meta.Valuer(user, ctx))
}
