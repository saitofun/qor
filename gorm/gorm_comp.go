package gorm

// @todo compatible gorm v1 and v2

// import (
// 	"reflect"
// 	"sync"
//
// 	gorm1 "github.com/jinzhu/gorm"
// )
//
// type Scope struct {
// 	Search          *search
// 	Value           interface{}
// 	SQL             string
// 	SQLVars         []interface{}
// 	db              *gorm1.DB
// 	instanceID      string
// 	primaryKeyField *gorm1.Field
// 	skipLeft        bool
// 	fields          *[]*gorm1.Field
// 	selectAttrs     *[]string
// }
//
// type search struct {
// 	db               *gorm1.DB
// 	whereConditions  []map[string]interface{}
// 	orConditions     []map[string]interface{}
// 	notConditions    []map[string]interface{}
// 	havingConditions []map[string]interface{}
// 	joinConditions   []map[string]interface{}
// 	initAttrs        []interface{}
// 	assignAttrs      []interface{}
// 	selects          map[string]interface{}
// 	omits            []string
// 	orders           []interface{}
// 	preload          []searchPreload
// 	offset           interface{}
// 	limit            interface{}
// 	group            string
// 	tableName        string
// 	raw              bool
// 	Unscoped         bool
// 	ignoreOrderQuery bool
// }
//
// type searchPreload struct {
// 	schema     string
// 	conditions []interface{}
// }
//
// type _Field struct {
// 	*StructField
// 	IsBlank bool
// 	Field   reflect.Value
// }
//
// type StructField struct {
// 	DBName          string
// 	Name            string
// 	Names           []string
// 	IsPrimaryKey    bool
// 	IsNormal        bool
// 	IsIgnored       bool
// 	IsScanner       bool
// 	HasDefaultValue bool
// 	Tag             reflect.StructTag
// 	TagSettings     map[string]string
// 	Struct          reflect.StructField
// 	IsForeignKey    bool
// 	Relationship    *Relationship
//
// 	tagSettingsLock sync.RWMutex
// }
//
// type Relationship struct {
// 	Kind                         string
// 	PolymorphicType              string
// 	PolymorphicDBName            string
// 	PolymorphicValue             string
// 	ForeignFieldNames            []string
// 	ForeignDBNames               []string
// 	AssociationForeignFieldNames []string
// 	AssociationForeignDBNames    []string
// 	JoinTableHandler             JoinTableHandlerInterface
// }
//
// type JoinTableHandlerInterface interface {
// 	// initialize join table handler
// 	Setup(relationship *Relationship, tableName string, source reflect.Type, destination reflect.Type)
// 	// Table return join table's table name
// 	Table(db *DB) string
// 	// Add create relationship in join table for source and destination
// 	Add(handler JoinTableHandlerInterface, db *DB, source interface{}, destination interface{}) error
// 	// Delete delete relationship in join table for sources
// 	Delete(handler JoinTableHandlerInterface, db *DB, sources ...interface{}) error
// 	// JoinWith query with `Join` conditions
// 	JoinWith(handler JoinTableHandlerInterface, db *DB, source interface{}) *DB
// 	// SourceForeignKeys return source foreign keys
// 	SourceForeignKeys() []JoinTableForeignKey
// 	// DestinationForeignKeys return destination foreign keys
// 	DestinationForeignKeys() []JoinTableForeignKey
// }
//
// // JoinTableForeignKey join table foreign key struct
// type JoinTableForeignKey struct {
// 	DBName            string
// 	AssociationDBName string
// }
//
