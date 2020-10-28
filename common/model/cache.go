package model

import (
	"go-api/common/global"
	"go-api/pkg/redis"
	"go-api/tools"
	"go-api/tools/config"
	"gorm.io/gorm"
	"reflect"
)

type Cache struct {
	CacheKey string `gorm:"-"`
}

func GetCacheRows(cachekey string) ([]interface{}, error) {
	var s []interface{}
	res, ok := global.Redis.Get(cachekey)
	if res != "" && ok == nil {
		s, _ = tools.JsonStrToSlice(res)
	}
	return s, nil
}

func (c *Cache) AfterFind(tx *gorm.DB) (err error) {
	// 1. 判断是否开缓存
	cacheSwitch := false
	cacheField := tx.Statement.Schema.LookUpField("CacheKey")
	if cacheField != nil && len(cacheField.BindNames) > 1 && cacheField.BindNames[0] == "Cache"{
		cacheSwitch = true
	}
	if !cacheSwitch{
		return nil
	}

	// 2. 缓存里面是否有
	sql := tx.Statement.SQL.String()
	ks := tx.Statement.ReflectValue.Kind().String()
	cachekey := redis.GetSqlKey(sql+ks)
	cacheRows, error := GetCacheRows(cachekey)
	if len(cacheRows) > 0 && error == nil {
		return nil
	}

	// 3. 判断主键cache设置expire秒数
	primaryField := tx.Statement.Schema.PrioritizedPrimaryField
	cacheTtl := primaryField.Tag.Get("cache")
	ttl := config.RedisConfig.Ttl
	if cacheTtl != "" {
		ttl, _ = tools.StringToInt(cacheTtl)
	}

	var s []interface{}
	var rj string
	switch tx.Statement.ReflectValue.Kind() {
	case reflect.Slice, reflect.Array:
		rj, _ = tools.SliceToJsonStr(tx.Statement.Dest)
	case reflect.Struct:
		s = append(s, tx.Statement.Dest)
		rj, _ = tools.SliceToJsonStr(s)
	}

	// 4. 插入redis
	error = global.Redis.Set(cachekey, rj, ttl)
	if error == nil {
		global.Logger.Info("db缓存插入成功", cachekey)
	}
	return
}
