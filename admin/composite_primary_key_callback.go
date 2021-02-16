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

		callbackProc := db.Callback().Query().Before("gorm:query")
		callbackName := "qor_admin:composite_primary_key"
		//if callbackProc.Get(callbackName) == nil {
		callbackProc.Register(callbackName, compositePrimaryKeyQueryCallback)
		//}

		callbackProc = db.Callback().Row().Before("gorm:row_query")
		// if callbackProc.Get(callbackName) == nil {
		callbackProc.Register(callbackName, compositePrimaryKeyQueryCallback)
		//}
	}
}

// DisableCompositePrimaryKeyMode disable composite primary key mode
var DisableCompositePrimaryKeyMode = "composite_primary_key:query:disable"

func compositePrimaryKeyQueryCallback(scope *gorm.DB) {
	if value, ok := scope.Get(DisableCompositePrimaryKeyMode); ok && value != "" {
		return
	}

	stmt := scope.Statement
	tableName := stmt.Table
	for _, primaryField := range stmt.Schema.PrimaryFields {
		if value, ok := scope.Get(fmt.Sprintf("primary_key[%v_%v]", tableName, primaryField.DBName)); ok && value != "" {
			stmt.Where(fmt.Sprintf("%v = ?", stmt.Quote(primaryField.DBName)), value)
			// scope.Search.Where(fmt.Sprintf("%v = ?", scope.Quote(primaryField.DBName)), value)
		}
	}
}
