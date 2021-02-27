package media

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/saitofun/qor/gorm"
	"github.com/saitofun/qor/serializable_meta"
)

var (
	// set MediaLibraryURL to change the default url /system/{{class}}/{{primary_key}}/{{column}}.{{extension}}
	MediaLibraryURL = ""
)

func cropField(field *gorm.Field, db *gorm.DB) (cropped bool) {
	fv, _ := field.ValueOf(reflect.ValueOf(db.Statement.Model))
	fval := reflect.ValueOf(fv)
	if fval.CanAddr() {
		// TODO Handle scanner
		if media, ok := fval.Addr().Interface().(Media); ok && !media.Cropped() {
			option := parseTagOption(field.Tag.Get("media_library"))
			if MediaLibraryURL != "" {
				option.Set("url", MediaLibraryURL)
			}
			if media.GetFileHeader() != nil || media.NeedCrop() {
				var mediaFile FileInterface
				var err error
				if fileHeader := media.GetFileHeader(); fileHeader != nil {
					mediaFile, err = media.GetFileHeader().Open()
				} else {
					mediaFile, err = media.Retrieve(media.URL("original"))
				}

				if err != nil {
					db.AddError(err)
					return false
				}

				media.Cropped(true)

				if url := media.GetURL(option, db, field, media); url == "" {
					db.AddError(errors.New("invalid URL"))
				} else {
					result, _ := json.Marshal(map[string]string{"Url": url})
					media.Scan(string(result))
				}

				if mediaFile != nil {
					defer mediaFile.Close()
					var handled = false
					for _, handler := range mediaHandlers {
						if handler.CouldHandle(media) {
							mediaFile.Seek(0, 0)
							if db.AddError(handler.Handle(media, mediaFile, option)) == nil {
								handled = true
							}
						}
					}

					// Save File
					if !handled {
						db.AddError(media.Store(media.URL(), option, mediaFile))
					}
				}
				return true
			}
		}
	}
	return false
}

func saveAndCropImage(isCreate bool) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		if db.Error == nil {
			var updateColumns = map[string]interface{}{}
			var dbVal = db.Statement.Model
			var dbValue = reflect.ValueOf(dbVal)

			// Handle SerializableMeta
			if value, ok := dbVal.(serializable_meta.SerializableMetaInterface); ok {
				var (
					isCropped        bool
					handleNestedCrop func(record interface{})
				)

				handleNestedCrop = func(record interface{}) {
					newSchema, _ := gorm.Parse(record)
					recordValue := reflect.ValueOf(record)
					for _, field := range newSchema.Fields {
						if cropField(field, db) {
							isCropped = true
							continue
						}
						fv, _ := field.ValueOf(recordValue)
						fVal := reflect.ValueOf(fv)

						if reflect.Indirect(fVal).Kind() == reflect.Struct {
							handleNestedCrop(fVal.Addr().Interface())
						}

						if reflect.Indirect(fVal).Kind() == reflect.Slice {
							for i := 0; i < reflect.Indirect(fVal).Len(); i++ {
								handleNestedCrop(reflect.Indirect(fVal).Index(i).Addr().Interface())
							}
						}
					}
				}

				record := value.GetSerializableArgument(value)
				handleNestedCrop(record)
				if isCreate && isCropped {
					updateColumns["value"], _ = json.Marshal(record)
				}
			}

			// Handle Normal Field
			for _, field := range db.Statement.Schema.Fields {
				if cropField(field, db) && isCreate {
					fv, _ := field.ValueOf(dbValue)
					updateColumns[field.DBName] = fv
				}
			}

			if db.Error == nil && len(updateColumns) != 0 {
				db.AddError(db.Session(&gorm.Session{}).Model(dbVal).UpdateColumns(updateColumns).Error)
			}
		}
	}
}

// RegisterCallbacks register callback into GORM DB
func RegisterCallbacks(db *gorm.DB) {
	if db.Callback().Create().Get("media:save_and_crop") == nil {
		db.Callback().Create().After("gorm:after_create").Register("media:save_and_crop", saveAndCropImage(true))
	}
	if db.Callback().Update().Get("media:save_and_crop") == nil {
		db.Callback().Update().Before("gorm:before_update").Register("media:save_and_crop", saveAndCropImage(false))
	}
}
