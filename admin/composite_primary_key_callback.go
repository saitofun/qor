package admin

import (
	"fmt"
	"regexp"

	"github.com/saitofun/qor/gorm"
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

		callbackName := "qor_admin:composite_primary_key"

		callbackProc := db.Callback().Query().Before("gorm:query")
		callbackProc.Register(callbackName, compositePrimaryKeyQueryCallback)

		callbackProc = db.Callback().Row().Before("gorm:row_query")
		callbackProc.Register(callbackName, compositePrimaryKeyQueryCallback)
	}
}

// DisableCompositePrimaryKeyMode disable composite primary key mode
var DisableCompositePrimaryKeyMode = "composite_primary_key:query:disable"

func compositePrimaryKeyQueryCallback(db *gorm.DB) {
	if v, ok := db.Get(DisableCompositePrimaryKeyMode); ok && v != "" {
		return
	}
	stmt := db.Statement
	schema := stmt.Schema
	tableName := stmt.Table
	for _, primaryField := range schema.PrimaryFields {
		v, ok := db.Get(fmt.Sprintf("primary_key[%v_%v]",
			tableName,
			primaryField.DBName))
		if ok && v != "" {
			stmt.BuildCondition(
				fmt.Sprintf("%v = ?", stmt.Quote(primaryField.DBName)), v)
		}
	}
}
