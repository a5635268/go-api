package database

import (
	"database/sql"
	"go-api/common/model"
	"go-api/pkg/redis"
	"gorm.io/gorm/callbacks"
	"log"
	"os"
	"reflect"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"go-api/common/config"
	"go-api/common/global"
	"go-api/tools"
	toolsConfig "go-api/tools/config"
)

type Mysql struct {
}

func (e *Mysql) Setup() {
	global.Source = e.GetConnect()
	global.Logger.Info(tools.Green(global.Source))
	db, err := sql.Open("mysql", global.Source)
	if err != nil {
		global.Logger.Fatal(tools.Red(e.GetDriver()+" connect error :"), err)
	}
	global.Cfg.SetDb(&config.DBConfig{
		Driver: "mysql",
		DB:     db,
	})
	global.Eloquent, err = e.Open(db, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		global.Logger.Fatal(tools.Red(e.GetDriver()+" connect error :"), err)
	} else {
		global.Logger.Info(tools.Green(e.GetDriver() + " connect success !"))
	}

	if global.Eloquent.Error != nil {
		global.Logger.Fatal(tools.Red(" database error :"), global.Eloquent.Error)
	}

	if toolsConfig.LoggerConfig.EnabledDB {
		global.Eloquent.Logger = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: time.Second,
			Colorful:      true,
			LogLevel:      logger.Info,
		})
	}

	// 注册全局回调
	global.Eloquent.Callback().Query().Replace("gorm:query", queryCallback)
}

// 打开数据库连接
func (e *Mysql) Open(db *sql.DB, cfg *gorm.Config) (*gorm.DB, error) {
	return gorm.Open(mysql.New(mysql.Config{Conn: db}), cfg)
}

// 获取数据库连接
func (e *Mysql) GetConnect() string {
	return toolsConfig.DatabaseConfig.Source
}

func (e *Mysql) GetDriver() string {
	return toolsConfig.DatabaseConfig.Driver
}


func queryCallback(db *gorm.DB) {
	if db.Error == nil {
		if db.Statement.Schema != nil && !db.Statement.Unscoped {
			for _, c := range db.Statement.Schema.QueryClauses {
				db.Statement.AddClause(c)
			}
		}

		if db.Statement.SQL.String() == "" {
			callbacks.BuildQuerySQL(db)
		}

		// 判断是否需要从缓存里面取
		cacheSwitch := false
		cacheField := db.Statement.Schema.LookUpField("CacheKey")
		if cacheField != nil && len(cacheField.BindNames) > 1 && cacheField.BindNames[0] == "Cache"{
			cacheSwitch = true
		}
		if cacheSwitch{
			sql := db.Statement.SQL.String()
			ks := db.Statement.ReflectValue.Kind().String()
			cachekey := redis.GetSqlKey(sql+ks)
			cacheRows,err := model.GetCacheRows(cachekey)
			if len(cacheRows) > 0 && err == nil{
				switch db.Statement.ReflectValue.Kind() {
				case reflect.Slice, reflect.Array:
					// 创建一个[]model的零值容器
					reflectValueType := db.Statement.ReflectValue.Type().Elem()
					elem := reflect.New(reflectValueType).Elem()
					db.Statement.ReflectValue.Set(reflect.MakeSlice(db.Statement.ReflectValue.Type(), 0, 0))

					// 开始填值
					for _, v := range cacheRows {
						m, _ := (v).(map[string]interface{})
						for _, field := range db.Statement.Schema.Fields {
							if len(field.BindNames) > 1 && field.BindNames[0] == "Cache"{
								continue
							}
							// 填充struct
							field.Set(elem, m[field.DBName])
						}
						// 填充slice
						db.Statement.ReflectValue.Set(reflect.Append(db.Statement.ReflectValue,elem))
					}

				case reflect.Struct:
					// 取第一个赋值
					m,_ := (cacheRows[0]).(map[string]interface{})
					for _, field := range db.Statement.Schema.Fields {
						if _, isZero := field.ValueOf(db.Statement.ReflectValue); isZero {
							// 判断一下，如果是来自cache的话，值就不用设置
							if len(field.BindNames) > 1 && field.BindNames[0] == "Cache"{
								continue
							}
							// 设置字段值
							field.Set(db.Statement.ReflectValue, m[field.DBName])
						}
					}
				}
				global.Logger.Info("数据结构从cache获取：", cacheRows)
				return
			}
		}

		// todo: 6. 缓存里面没有，数据库里数据写进缓存
		if !db.DryRun && db.Error == nil {
			global.Logger.Info("从db获取：")
			rows, err := db.Statement.ConnPool.QueryContext(db.Statement.Context, db.Statement.SQL.String(), db.Statement.Vars...)
			if err != nil {
				db.AddError(err)
				return
			}
			defer rows.Close()

			gorm.Scan(rows, db, false)
		}
	}
}
