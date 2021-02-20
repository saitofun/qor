package resource

import (
	"errors"
	"fmt"
	"reflect"
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
func (res *Resource) ToPrimaryQueryParams(value string, ctx *qor.Context) (query string, args []interface{}) {
	if value == "" {
		return
	}

	var (
		stmt      = ctx.GetDB().Session(&gorm.Session{}).Statement
		fArgs     = strings.Split(value, ",")
		schema, _ = gorm.Parse(res.Value)
	)

	if len(fArgs) == len(res.PrimaryFields) {
		var cond []string
		for i, f := range res.PrimaryFields {
			cond = append(cond, fmt.Sprintf("%v.%v = ?",
				stmt.Quote(schema.Table), stmt.Quote(f.DBName)))
			args = append(args, fArgs[i])
		}
		query = strings.Join(cond, " AND ")
		return

	} else if f := res.primaryField; f != nil {
		return fmt.Sprintf("%v.%v = ?",
			stmt.Quote(stmt.Table), stmt.Quote(f.DBName)), []interface{}{value}
	}

	return
}

// ToPrimaryQueryParamsFromMetaValue generate query params based on MetaValues
func (res *Resource) ToPrimaryQueryParamsFromMetaValue(metas *MetaValues, ctx *qor.Context) (query string, args []interface{}) {
	var cond []string
	if metas == nil {
		return
	}

	stmt := ctx.GetDB().Statement
	schema, _ := gorm.Parse(res.Value)
	for _, f := range res.PrimaryFields {
		if mf := metas.Get(f.Name); mf != nil {
			cond = append(cond, fmt.Sprintf("%v.%v = ?",
				stmt.Quote(schema.Table), stmt.Quote(f.DBName)))
			args = append(args, utils.ToString(mf.Value))
		}
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
	if res.HasPermission(roles.Read, context) {
		db := context.GetDB()
		if _, ok := db.Get("qor:getting_total_count"); ok {
			_result := new(int64)
			err := context.GetDB().Count(_result).Error
			switch result.(type) {
			case *int:
				*(result).(*int) = int(*_result)
			case *int64:
				*(result).(*int64) = *_result
			}
			return err
		}
		return context.GetDB().Set("gorm:order_by_primary_key", "DESC").Find(result).Error
	}

	return roles.ErrPermissionDenied
}

func (res *Resource) saveHandler(result interface{}, ctx *qor.Context) error {
	schema, _ := gorm.Parse(res.Value)
	_, zero := schema.PrioritizedPrimaryField.ValueOf(reflect.ValueOf(res.Value))

	if (zero &&
		res.HasPermission(roles.Create, ctx)) || // has create permission
		res.HasPermission(roles.Update, ctx) { // has update permission
		return ctx.GetDB().Save(result).Error
	}
	return roles.ErrPermissionDenied
}

func (res *Resource) deleteHandler(result interface{}, context *qor.Context) error {
	if res.HasPermission(roles.Delete, context) {
		if sql, args := res.ToPrimaryQueryParams(context.ResourceID, context); sql != "" {
			db := context.GetDB().Session(&gorm.Session{})
			if !errors.Is(db.First(result, append([]interface{}{sql}, args...)...).Error, gorm.ErrRecordNotFound) {
				return db.Delete(result).Error
			}
		}
		return gorm.ErrRecordNotFound
	}
	return roles.ErrPermissionDenied
}
