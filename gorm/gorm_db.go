package gorm

import (
	"bytes"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type DB = gorm.DB

var Open = gorm.Open

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

var commonInitialisms = []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "TTL", "UID", "UI", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}
var commonInitialismsReplacer *strings.Replacer

func init() {
	var commonInitialismsForReplacer []string
	for _, initialism := range commonInitialisms {
		commonInitialismsForReplacer = append(commonInitialismsForReplacer, initialism, strings.Title(strings.ToLower(initialism)))
	}
	commonInitialismsReplacer = strings.NewReplacer(commonInitialismsForReplacer...)
}

func ToDBName(name string) string {
	const (
		lower = false
		upper = true
	)
	var (
		value      = commonInitialismsReplacer.Replace(name)
		buf        = bytes.NewBufferString("")
		lastCase   bool
		currCase   bool
		nextCase   bool
		nextNumber bool
	)

	for i, v := range value[:len(value)-1] {
		nextCase = value[i+1] >= 'A' && value[i+1] <= 'Z'
		nextNumber = value[i+1] >= '0' && value[i+1] <= '9'

		if i > 0 {
			if currCase == upper {
				if lastCase == upper && (nextCase == upper || nextNumber == upper) {
					buf.WriteRune(v)
				} else {
					if value[i-1] != '_' && value[i+1] != '_' {
						buf.WriteRune('_')
					}
					buf.WriteRune(v)
				}
			} else {
				buf.WriteRune(v)
				if i == len(value)-2 && (nextCase == upper && nextNumber == lower) {
					buf.WriteRune('_')
				}
			}
		} else {
			currCase = upper
			buf.WriteRune(v)
		}
		lastCase = currCase
		currCase = nextCase
	}

	buf.WriteByte(value[len(value)-1])

	return strings.ToLower(buf.String())
}
