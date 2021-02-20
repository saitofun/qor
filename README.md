# gorm迁移至v2版本改动说明

1. 错误判断 error check

v1中所有错误类型已经在v2中重命名
如：
```go
RecordNotFound => ErrRecordNotFound
```

v2错误判断在示例
```go
err := db.First(&user).Error
errors.Is(err, gorm.ErrRecordNotFound)
```


2. 模型解析 modle parse

v1:
```
scope := &grom.Scope{Value:model}
```
v2:
```
// 缺省使用默认的命名策略
schema, _ = schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
```
schema解析后会持有model对应数据库模型的所有信息，包括主键，表名，反射类型，关联关系等

3. 关联查询 assosiation query

v1:
```
db.Model(model).Related(subModel, subFieldName) // v2中已经废弃不用了
```
v2:
```
db.Model(model).Association(subFieldName).Find(subMode)
```




4. 关联关系确定 relationships

使用schema.Parse解析成功后， schema.Relationships.Relations中会包含对应field.DBName与被解析模型的关联关系

关联类型与v1基本保持不变，包括HasOne，HasMany，BelongsTo，Many2Many四种

对于多对多的关联需要join的时候 schema.Relationships.Relations[FieldDBName].JoinTable中会持有被关联模型相关上下文

5. 数据库迁移 migration

v2版本中的Migrator接口不包含AddUniqueIndex方法。

如果需要在程序执行中创建唯一索引，目前建议使用RawSQL

6. 子句 clause

v2版本中对子句更加严谨，并有了独立封装, 在gorm.io/clause 中

包含curd，条件， limit ， 排序等。。。

v2继续支持使用raw sql，但不推荐


7. 钩子

v1回调接口： 
```go
func(*gorm.Scope) {}
```

v2回调接口
```go
func(*gorm.DB) {}
```

v2中默认回调注册在 `callbacks.RegisterDefaultCallback` 方法中

qor/validations 中对字段的预校验方法不推荐使用

建议使用gorm的check tag


8. 方言

v2版本中方言全部迁移至`gorm.io/gorm/driver/...`

gorm.Open 与v1版本不再兼容, 并返回`*gorm.DB`，而不是`gorm.DB`，

需使用Dialector接口与gorm.Config作为参数

```go
import "gorm.io/driver/sqlite"

dsn:= "path_to_sqlite_db_file"
db := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
```

9 其他

如果需要维持db的上下文不受影响， 使用session，保证线程安全
```go
db := db.Session(&gorm.Session{})
```

不建议使用db.Statement.Parse, 该方法会使用schema.Parse方法，但会传入db上下文
```go
db.Statement.Parse(model) // 不推荐
schema.Parse(model, &sync.Map{}, schema.Namer) // 推荐
```
