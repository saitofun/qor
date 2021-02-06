package qor

import "gorm.io/gorm"

// Config qor config struct
type Config struct {
	DB *gorm.DB
}
