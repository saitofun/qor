package resource

import (
	"errors"
	"fmt"
	"strings"

	"github.com/saitofun/qor/gorm"
	"github.com/saitofun/qor/qor"
	"github.com/saitofun/qor/qor/utils"
	"github.com/saitofun/qor/roles"
)

// CallFindOne call find one method
func (res *Resource) CallFindOne(result interface{}, metaValues *MetaValues, context *qor.Context) error {
	return res.FindOneHandler(result, metaValues, context)
}

// CallFindMany call find many method
func (res *Resource) CallFindMany(result interface{}, context *qor.Context) error {
	return res.FindManyHandler(result, context)
}

// CallSave call save method
func (res *Resource) CallSave(result interface{}, context *qor.Context) error {
	return res.SaveHandler(result, context)
}

// CallDelete call delete method
func (res *Resource) CallDelete(result interface{}, context *qor.Context) error {
	return res.DeleteHandler(result, context)
}

// ToPrimaryQueryParams generate query params based on primary key, multiple primary value are linked with a comma
func (res *Resource) ToPrimaryQueryParams(primaryValue string, context *qor.Context) (sql string, primaries []interface{}) {
	if primaryValue == "" {
		return
	}

	var (
		stmt    = context.DB.Statement
		table   = stmt.Table
		clauses []string
	)
	quoter := func(f *gorm.Field) {
		clauses = append(clauses, fmt.Sprintf("%v.%v = ?", stmt.Quote(table), stmt.Quote(f.DBName)))
		primaries = append(primaries, gorm.ReflectFieldValue(res.Value, f))
	}

	if len(res.PrimaryFields) > 1 {
		primaryStrings := strings.Split(primaryValue, ",")
		if len(primaryStrings) == len(res.PrimaryFields) {
			for _, f := range res.PrimaryFields {
				quoter(f)
			}
		}
	} else if f := res.primaryField; f != nil {
		quoter(f)
	} else {
		schema, _ := gorm.ModelToSchema(res.Value)
		quoter(schema.PrioritizedPrimaryField)
	}

	if len(clauses) > 0 {
		sql = strings.Join(clauses, " AND ")
	}

	return
}

// ToPrimaryQueryParamsFromMetaValue generate query params based on MetaValues
func (res *Resource) ToPrimaryQueryParamsFromMetaValue(metaValues *MetaValues, context *qor.Context) (sql string, fields []interface{}) {
	var (
		stmt    = context.DB.Statement
		table   = stmt.Table
		clauses []string
	)

	quoter := func(f *gorm.Field) {
		clauses = append(clauses, fmt.Sprintf("%v.%v = ?", stmt.Quote(table), stmt.Quote(f.DBName)))
		fields = append(fields, utils.ToString(metaValues.Get(f.Name).Value))
	}

	if metaValues == nil {
		return
	}

	for _, field := range res.PrimaryFields {
		if m := metaValues.Get(field.Name); m != nil {
			quoter(field)
		}
	}

	if len(clauses) > 0 {
		sql = strings.Join(clauses, " AND ")
	}

	return
}

func (res *Resource) findOneHandler(result interface{}, metaValues *MetaValues, context *qor.Context) error {
	if res.HasPermission(roles.Read, context) {
		var (
			primaryQuerySQL string
			primaryParams   []interface{}
		)

		if metaValues == nil {
			primaryQuerySQL, primaryParams = res.ToPrimaryQueryParams(context.ResourceID, context)
		} else {
			primaryQuerySQL, primaryParams = res.ToPrimaryQueryParamsFromMetaValue(metaValues, context)
		}

		if primaryQuerySQL != "" {
			if metaValues != nil {
				if destroy := metaValues.Get("_destroy"); destroy != nil {
					if fmt.Sprint(destroy.Value) != "0" && res.HasPermission(roles.Delete, context) {
						context.GetDB().Delete(result, append([]interface{}{primaryQuerySQL}, primaryParams...)...)
						return ErrProcessorSkipLeft
					}
				}
			}
			return context.GetDB().First(result, append([]interface{}{primaryQuerySQL}, primaryParams...)...).Error
		}

		return errors.New("failed to find")
	}
	return roles.ErrPermissionDenied
}

func (res *Resource) findManyHandler(result interface{}, context *qor.Context) error {
	if !res.HasPermission(roles.Read, context) {
		return roles.ErrPermissionDenied
	}
	v, ok := result.(*int64)
	if !ok || v == nil {
		return errors.New("unexpected parameter: result should be *int64")
	}
	if _, ok := context.GetDB().Get("qor:getting_total_count"); ok {
		return context.GetDB().Count(v).Error
	}
	return context.GetDB().Set("gorm:order_by_primary_key", "DESC").Find(result).Error
}

func (res *Resource) saveHandler(result interface{}, context *qor.Context) error {
	schema, _ := gorm.ModelToSchema(result)
	if schema.PrioritizedPrimaryField == nil &&
		res.HasPermission(roles.Create, context) ||
		res.HasPermission(roles.Update, context) {
		return context.DB.Save(result).Error
	}
	return roles.ErrPermissionDenied
}

func (res *Resource) deleteHandler(result interface{}, context *qor.Context) error {
	if res.HasPermission(roles.Delete, context) {
		sql, values := res.ToPrimaryQueryParams(context.ResourceID, context)
		if sql == "" {
			return gorm.ErrRecordNotFound
		}
		err := context.DB.First(result,
			append([]interface{}{sql}, values...)...).Error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return context.DB.Delete(result).Error
		}
	}
	return gorm.ErrRecordNotFound
}
