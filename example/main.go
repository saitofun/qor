package main

import (
	"fmt"
	"os"

	"github.com/kataras/iris/v12"
	_ "github.com/mattn/go-sqlite3"
	"github.com/saitofun/qor/admin"
	"github.com/saitofun/qor/gorm"
	"gorm.io/driver/sqlite"
)

// Create a GORM-backend model
type User struct {
	gorm.Model
	Name string
}

// Create another GORM-backend model
type Product struct {
	gorm.Model
	Name        string
	Description string
}

func main() {
	fmt.Println(os.Getwd())

	DB, _ := gorm.Open(sqlite.Open("demo.db"), &gorm.Config{})
	DB.AutoMigrate(&User{}, &Product{})

	qorPrefix := "/admin"
	// Initialize Qor Admin.
	Admin := admin.New(&admin.AdminConfig{DB: DB})

	// Allow to use Admin to manage User, Product
	Admin.AddResource(&User{})
	Admin.AddResource(&Product{})
	// Create a qor handler and convert it to an iris one with `iris.FromStd`.
	handler := iris.FromStd(Admin.NewServeMux(qorPrefix))

	// Initialize Iris.
	app := iris.New()
	// Mount routes for "/admin" and "/admin/:xxx/..."
	app.Any(qorPrefix, handler)
	app.Any(qorPrefix+"/{p:path}", handler)

	// Start the server.
	// Navigate at: http://localhost:9000/admin.
	app.Listen(":9000")
}
