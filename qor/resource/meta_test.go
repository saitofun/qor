package resource_test

import (
	"database/sql/driver"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/saitofun/qor/gorm"
	"github.com/saitofun/qor/qor"
	"github.com/saitofun/qor/qor/resource"
	"github.com/saitofun/qor/qor/utils"
	"github.com/saitofun/qor/utils/test_db"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

func format(value interface{}) string {
	return fmt.Sprint(utils.Indirect(reflect.ValueOf(value)).Interface())
}

func checkMeta(record interface{}, meta *resource.Meta, value interface{}, t *testing.T, expectedValues ...string) {
	var (
		context       = &qor.Context{DB: test_db.NewTestDB()}
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

var DB *gorm.DB

func init() {
	var err error
	if DB, err = OpenTestConnection(); err != nil {
		log.Printf("failed to connect database, got error %v", err)
		os.Exit(1)
	} else {
		sqlDB, err := DB.DB()
		if err == nil {
			err = sqlDB.Ping()
		}

		if err != nil {
			log.Printf("failed to connect database, got error %v", err)
		}

		RunMigrations()
		if DB.Dialector.Name() == "sqlite" {
			DB.Exec("PRAGMA foreign_keys = ON")
		}
	}
}

func OpenTestConnection() (db *gorm.DB, err error) {
	cfg := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	// log.Println("testing mysql...")
	// db, err = gorm.Open(mysql.Open("root:root@tcp(localhost:33060)/gorm_test?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})

	log.Println("testing sqlite3...")
	db, err = gorm.Open(sqlite.Open("/Users/sincos/sincos/src/github.com/saitofun/gorm.db"), cfg)

	db.Logger = db.Logger.LogMode(logger.Info)

	return
}

func RunMigrations() {
	var err error
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(models), func(i, j int) { models[i], models[j] = models[j], models[i] })

	DB.Migrator().DropTable("user_friends", "user_speaks")

	if err = DB.Migrator().DropTable(models...); err != nil {
		log.Printf("Failed to drop table, got error %v\n", err)
		os.Exit(1)
	}

	if err = DB.AutoMigrate(models...); err != nil {
		log.Printf("Failed to auto migrate, but got error %v\n", err)
		os.Exit(1)
	}

	for _, m := range models {
		if !DB.Migrator().HasTable(m) {
			log.Printf("Failed to create table for %#v\n", m)
			os.Exit(1)
		}
	}
}

type CreditCard struct {
	gorm.Model
	Number string
	Issuer string
}

type Company struct {
	gorm.Model
	Name string
}

type Address struct {
	gorm.Model
	UserID   uint
	Address1 string
	Address2 string
}

type Language struct {
	gorm.Model
	Name string
}

type User struct {
	gorm.Model
	Name         string `gorm:"size:50"`
	Age          uint
	Role         string
	Active       bool
	RegisteredAt *time.Time
	Profile      Profile    // has one
	CreditCardID uint       // foreign key
	CreditCard   CreditCard // belongs to
	Addresses    []Address  // has many
	CompanyID    uint       // foreign key
	Company      *Company   // belongs to
	Languages    []Language `gorm:"many2many:user_languages;"` // many 2 many
}

type Profile struct {
	gorm.Model
	UserID uint
	Name   string
	Sex    string

	Phone Phone
}

type Phone struct {
	gorm.Model

	ProfileID uint64
	Num       string
}

var models = []interface{}{
	&User{},
	&CreditCard{},
	&Company{},
	&Address{},
	&Language{},
	&Profile{},
	&Phone{},
}

func TestMeta_PreInitialize(t *testing.T) {
	var (
		model = &User{}
		res   = resource.New(model)
		meta  *resource.Meta
	)
	meta = &resource.Meta{
		Name:         "age", // this is not a field name
		FieldName:    "age",
		BaseResource: res,
	}
	err := meta.PreInitialize()
	if err == nil || err.Error() != "meta resource_test.User no field: age" {
		t.Errorf("error should catch because not a field name")
		return
	}

	meta.Name, meta.FieldName = "Age", "Age"
	err = meta.PreInitialize()
	if err != nil {
		t.Errorf("error should catch nil, but %v", err)
		return
	}

	meta.Name, meta.FieldName = "Company.ID", "Company.ID"
	err = meta.PreInitialize()
	if err != nil {
		t.Errorf("error should catch nil, but %v", err)
		return
	}
	f := meta.FieldStruct
	if f.Schema.Name != "User" {
		t.Errorf("f.Schema.Name should be `User`, but `%s`", f.Schema.Name)
		return
	}
	if f.Name != "ID" {
		t.Errorf("f.Name should be `ID`, but `%s`", f.Name)
		return
	}
	if f.DBName != "id" {
		t.Errorf("f.DBName should be `id`, but `%s`", f.DBName)
		return
	}

	meta.Name, meta.FieldName = "company", "Company"
	err = meta.PreInitialize()
	if err != nil {
		t.Errorf("error should catch nil, but %v", err)
		return
	}
	f = meta.FieldStruct
	if f.Schema.Name != "User" {
		t.Errorf("f.Schema.Name should be `User`, but `%s`", f.Schema.Name)
		return
	}
	if f.Name != "Company" {
		t.Errorf("f.Name should be `Company`, but `%s`", f.Name)
		return
	}
	if f.DBName != "" {
		t.Errorf("f.DBName should be `\"\"`, but `%s`", f.DBName)
		return
	}
	if reflect.TypeOf((*Company)(nil)) != f.FieldType {
		t.Errorf("f.FieldType should be `*resource_test.Company`, but `%v`", f.FieldType)
		return
	}
	if reflect.TypeOf(Company{}) != f.IndirectFieldType {
		t.Errorf("f.IndirectFieldType should be `resource_test.Company`, but `%v`", f.IndirectFieldType)
		return
	}
	t.Log(f.Schema.Relationships.Relations[f.Name].Name)
	t.Log(f.Schema.Relationships.Relations[f.Name].Type)
	t.Log(f.Schema.Relationships.Relations[f.Name].Schema.Name)
	t.Log(f.Schema.Relationships.Relations[f.Name].Schema.Table)
	t.Log(f.Schema.Relationships.Relations[f.Name].FieldSchema.Name)
	t.Log(f.Schema.Relationships.Relations[f.Name].FieldSchema.Table)
}

func TestMeta_ValuerForDirectField(t *testing.T) {
	var (
		user = &User{Age: 18}
		res  = resource.New(user)
		meta *resource.Meta
		ctx  = &qor.Context{DB: DB}
	)
	meta = &resource.Meta{
		Name:         "Age", // this is not a field name
		FieldName:    "Age",
		BaseResource: res,
	}
	meta.PreInitialize()
	meta.Initialize()

	result := meta.Valuer(user, ctx)
	if v, ok := result.(uint); !ok || v != 18 {
		t.Errorf("should got user's Age is 18, but got %v", result)
		return
	}

	user.Name = "sai"
	DB.Save(user)
	DB.Find(user, clause.Eq{
		Column: "name",
		Value:  "sai",
	})
	if user.Name != "sai" || user.Age != 18 {
		t.Errorf("db save error: %v", DB.Error)
		return
	}
	user.Age, user.Name = 0, ""
	result = meta.Valuer(user, ctx)
	if v, ok := result.(uint); !ok || v != 18 {
		t.Errorf("should got user's Age is 18, but got %v", result)
		return
	}
}

func TestMeta_ValuerForNestedField_UserProfile_HasOne(t *testing.T) {
	var (
		user = &User{}
		res  = resource.New(user)
		meta *resource.Meta
		ctx  = &qor.Context{DB: DB}
		db   = DB
	)

	user = &User{
		Name:   "sai",
		Age:    18,
		Role:   "manager",
		Active: true,
	}
	db.Save(user)
	DB.Find(user, clause.Eq{
		Column: "name",
		Value:  "sai",
	})
	if user.Name != "sai" || user.Age != 18 {
		t.Errorf("db save error: %v", DB.Error)
		return
	}
	profile := &Profile{
		UserID: user.ID,
		Name:   "sai_profile",
		Sex:    "M",
		Phone:  Phone{Num: "18888888888"},
	}
	db.Save(profile)

	meta = &resource.Meta{
		Name:         "user's Phone Number",
		FieldName:    "Profile.Phone.Num",
		BaseResource: res,
	}
	meta.PreInitialize()
	meta.Initialize()

	result := meta.Valuer(user, ctx)
	if v, ok := result.(string); !ok || v != "18888888888" {
		t.Errorf("user's number should be `18888888888` but got %v", v)
		return
	}

}

func TestMeta_ValuerForNestedField_CreditCard_BelongsTo(t *testing.T) {
	var (
		user = &User{}
		res  = resource.New(user)
		meta *resource.Meta
		ctx  = &qor.Context{DB: DB}
		db   = DB
	)
	// test CreditCard.number
	user = &User{
		Name:   "sai",
		Age:    18,
		Role:   "manager",
		Active: true,
		CreditCard: CreditCard{
			Number: "8888 8888 8888 8887",
			Issuer: "CMBC",
		},
	}
	db.Save(user)
	user2 := &User{}
	db.Find(user2, clause.Eq{
		Column: "ID",
		Value:  user.ID,
	})

	meta = &resource.Meta{
		Name:         "user's credit card",
		FieldName:    "CreditCard.Number",
		BaseResource: res,
	}

	meta.PreInitialize()
	meta.Initialize()

	result := meta.Valuer(&User{Model: gorm.Model{ID: user2.ID}}, ctx)
	if v, ok := result.(string); !ok || v != "8888 8888 8888 8887" {
		t.Errorf("user's credit card number should be `\"8888 8888 8888 8887\"`,"+
			" but got : %v", v)
		return
	}
}

func TestMeta_ValuerForNestedPtrField_Company_BelongsTo(t *testing.T) {
	var (
		user = &User{}
		res  = resource.New(user)
		meta *resource.Meta
		ctx  = &qor.Context{DB: DB}
		db   = DB
	)
	// test CreditCard.number
	user = &User{
		Name:   "sai",
		Age:    18,
		Role:   "manager",
		Active: true,
		Company: &Company{
			Name: "SAI",
		},
	}
	db.Save(user)
	user2 := &User{}
	db.Find(user2, clause.Eq{
		Column: "ID",
		Value:  user.ID,
	})

	meta = &resource.Meta{
		Name:         "user's company name",
		FieldName:    "Company.Name",
		BaseResource: res,
	}

	meta.PreInitialize()
	meta.Initialize()

	result := meta.Valuer(&User{Model: gorm.Model{ID: user2.ID}}, ctx)
	if v, ok := result.(string); !ok || v != "SAI" {
		t.Errorf("user's company name should be `\"SAI\"`, but got : %v", v)
		return
	}
}
