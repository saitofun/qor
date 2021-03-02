package admin

import (
	"fmt"
	"regexp"

	"github.com/saitofun/qor/gorm"
	"gorm.io/gorm/clause"
)

var primaryKeyRegexp = regexp.MustCompile(`primary_key\[.+_.+\]`)

func (admin Admin) registerCompositePrimaryKeyCallback() {
	if db := admin.DB; db != nil {
		// register middleware
		router := admin.GetRouter()
		router.Use(&Middleware{
			Name: "composite primary key filter",
			Handler: func(context *Context, middleware *Middleware) {
				db := context.GetDB()
				for key, value := range context.Request.URL.Query() {
					if primaryKeyRegexp.MatchString(key) {
						db = db.Set(key, value)
					}
				}
				context.SetDB(db)

				middleware.Next(context)
			},
		})

		name := "qor_admin:composite_primary_key"

		proc := db.Callback().Query().Before("gorm:query")
		proc.Register(name, compositePrimaryKeyQueryCallback)

		proc = db.Callback().Row().Before("gorm:row_query")
		proc.Register(name, compositePrimaryKeyQueryCallback)
	}
}

// DisableCompositePrimaryKeyMode disable composite primary key mode
var DisableCompositePrimaryKeyMode = "composite_primary_key:query:disable"

func compositePrimaryKeyQueryCallback(db *gorm.DB) {
	value, ok := db.Get(DisableCompositePrimaryKeyMode)
	if value == nil || ok && value != "" {
		return
	}
	stmt := db.Statement
	schema, _ := gorm.Parse(value)
	for _, pf := range schema.PrimaryFields {
		pk := fmt.Sprintf("[primary_key[%v_%v]", stmt.Table, pf.DBName)
		if v, ok := db.Get(pk); ok && v != "" {
			db = db.Where(clause.Eq{Column: pf.DBName, Value: v})
		}
	}
}
