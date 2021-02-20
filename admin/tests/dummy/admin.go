package dummy

import (
	"errors"
	"fmt"

	"github.com/saitofun/qor/admin"
	"github.com/saitofun/qor/gorm"
	"github.com/saitofun/qor/media"
	"github.com/saitofun/qor/qor"
	"github.com/saitofun/qor/utils/test_db"
)

// NewDummyAdmin generate admin for dummy app
func NewDummyAdmin(keepData ...bool) *admin.Admin {
	var (
		db     = test_db.NewTestDB()
		models = []interface{}{&User{}, &CreditCard{}, &Address{}, &Language{}, &Profile{}, &Phone{}, &Company{}}
		Admin  = admin.New(&qor.Config{DB: db})
	)

	media.RegisterCallbacks(db)

	for _, value := range models {
		if len(keepData) == 0 {
			_ = db.Migrator().DropTable(value)
		}
		db.AutoMigrate(value)
	}

	Admin.AddResource(&Company{})
	Admin.AddResource(&Language{}, &admin.Config{Name: "语种 & 语言", Priority: -1})
	user := Admin.AddResource(&User{})
	user.Meta(&admin.Meta{
		Name: "CreditCard",
		Type: "single_edit",
	})
	user.Meta(&admin.Meta{
		Name: "Languages",
		Type: "select_many",
		Collection: func(resource interface{}, context *qor.Context) (results [][]string) {
			var (
				languages []Language
				err       = context.GetDB().Find(&languages).Error
			)
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				for _, language := range languages {
					results = append(results, []string{fmt.Sprint(language.ID), language.Name})
				}
			}
			return
		},
	})

	return Admin
}
