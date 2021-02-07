package validations_test

import (
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/asaskevich/govalidator"
	_ "github.com/mattn/go-sqlite3"
	"github.com/saitofun/qor/gorm"
	"github.com/saitofun/qor/utils/test_utils"
	"github.com/saitofun/qor/validations"
)

var db *gorm.DB

type User struct {
	gorm.Model
	Name           string `valid:"required"`
	Password       string `valid:"length(6|20)"`
	SecurePassword string `valid:"numeric"`
	Email          string `valid:"email~Email already be token"`
	CompanyID      int
	Company        Company
	CreditCard     CreditCard
	Addresses      []Address
	Languages      []Language `gorm:"many2many:user_languages"`
}

func (user *User) Validate(db *gorm.DB) {
	if user.Name == "invalid" {
		db.AddError(validations.NewError(user, "Name", "invalid user name"))
	}
	govalidator.CustomTypeTagMap.Set("email",
		func(email interface{}, context interface{}) bool {
			if email.(string) == "" {
				return true
			}
			var count int64
			db.Model(&User{}).Where("email = ?", email).Count(&count)
			if count == 0 {
				return true
			}
			return false
		},
	)
}

type Company struct {
	gorm.Model
	Name string
}

func (company *Company) Validate(db *gorm.DB) {
	if company.Name == "invalid" {
		db.AddError(errors.New("invalid company name"))
	}
}

type CreditCard struct {
	gorm.Model
	UserID int
	Number string
}

func (card *CreditCard) Validate(db *gorm.DB) {
	if !regexp.MustCompile("^(\\d){13,16}$").MatchString(card.Number) {
		db.AddError(validations.NewError(card, "Number", "invalid card number"))
	}
}

type Address struct {
	gorm.Model
	UserID  int
	Address string
}

func (address *Address) Validate(db *gorm.DB) {
	if address.Address == "invalid" {
		db.AddError(validations.NewError(address, "Address", "invalid address"))
	}
}

type Language struct {
	gorm.Model
	Code string
}

func (language *Language) Validate(db *gorm.DB) error {
	if language.Code == "invalid" {
		return validations.NewError(language, "Code", "invalid language")
	}
	return nil
}

func init() {
	db = test_utils.TestDB()
	validations.RegisterCallbacks(db)
	tabs := []interface{}{
		&User{},
		&Company{},
		&CreditCard{},
		&Address{},
		&Language{},
	}
	for _, t := range tabs {
		if e := db.Migrator().DropTable(t); e != nil {
			panic(e)
		}
		if e := db.AutoMigrate(t); e != nil {
			panic(e)
		}
	}
}

func TestGoValidation(t *testing.T) {
	// defer db.Delete(&User{})

	user := User{Name: "", Password: "123123", Email: "a@gmail.com"}

	result := db.Save(&user)
	if result.Error == nil {
		t.Errorf("Should get error when save empty user")
	}

	errs := strings.Split(result.Error.Error(), ";")

	if len(errs) <= 0 && errs[0] != "Name can't be blank" {
		t.Errorf("Error message should be equal `Name can't be blank`")
	}

	user = User{Name: "", Password: "123", SecurePassword: "AB123", Email: "aagmail.com"}
	result = db.Save(&user)
	messages := []string{
		"Name can't be blank",
		"Password is the wrong length (should be 6~20 characters)",
		"SecurePassword is not a number",
		"Email is not a valid email address",
	}
	errs = strings.Split(result.Error.Error(), ";")
	if len(errs) != len(messages) {
		t.Errorf("unexpected error count")
	}
	for i, err := range errs {
		if messages[i] != strings.TrimSpace(err) {
			t.Errorf("Error message should be equal `%v`, but it is `%v`",
				messages[i], err)
		}
	}

	user = User{Name: "A", Password: "123123", Email: "a@gmail.com"}
	ret := db.Save(&user)
	if ret.Error != nil {
		t.Error(ret.Error)
	}

	user = User{Name: "B", Password: "123123", Email: "a@gmail.com"}
	if err := db.Save(&user).Error; err == nil || err.Error() != "Email already be token" {
		t.Errorf("Should get email alredy be token error")
	}
}

func TestSaveInvalidUser(t *testing.T) {
	user := User{Name: "invalid"}

	if result := db.Save(&user); result.Error == nil {
		t.Errorf("Should get error when save invalid user")
	}
}

func TestSaveInvalidCompany(t *testing.T) {
	user := User{
		Name:    "valid",
		Company: Company{Name: "invalid"},
	}

	if result := db.Save(&user); result.Error == nil {
		t.Errorf("Should get error when save invalid company")
	}
}

func TestSaveInvalidCreditCard(t *testing.T) {
	user := User{
		Name:       "valid",
		Company:    Company{Name: "valid"},
		CreditCard: CreditCard{Number: "invalid"},
	}

	if result := db.Save(&user); result.Error == nil {
		t.Errorf("Should get error when save invalid credit card")
	}
}

func TestSaveInvalidAddresses(t *testing.T) {
	user := User{
		Name:       "valid",
		Company:    Company{Name: "valid"},
		CreditCard: CreditCard{Number: "4111111111111111"},
		Addresses:  []Address{{Address: "invalid"}},
	}

	if result := db.Save(&user); result.Error == nil {
		t.Errorf("Should get error when save invalid addresses")
	}
}

func TestSaveInvalidLanguage(t *testing.T) {
	user := User{
		Name:       "valid",
		Company:    Company{Name: "valid"},
		CreditCard: CreditCard{Number: "4111111111111111"},
		Addresses:  []Address{{Address: "valid"}},
		Languages:  []Language{{Code: "invalid"}},
	}

	if result := db.Save(&user); result.Error == nil {
		t.Errorf("Should get error when save invalid language")
	}
}

func TestSaveAllValidData(t *testing.T) {
	user := User{
		Name:       "valid",
		Company:    Company{Name: "valid"},
		CreditCard: CreditCard{Number: "4111111111111111"},
		Addresses:  []Address{{Address: "valid1"}, {Address: "valid2"}},
		Languages:  []Language{{Code: "valid1"}, {Code: "valid2"}},
	}

	if result := db.Save(&user); result.Error != nil {
		t.Errorf("Should get no error when save valid data, but got: %v", result.Error)
	}
}
