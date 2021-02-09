package gorm

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type DB = gorm.DB

func GetDBErrors(db *DB) (ret []error) {
	if db.Error == nil {
		return nil
	}
	msg := strings.Split(db.Error.Error(), ";")
	for _, v := range msg {
		ret = append(ret, errors.New(strings.TrimSpace(v)))
	}
	return
}
