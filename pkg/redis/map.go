package redis

import (
	"crypto/md5"
	"fmt"
	"io"
)

const (
	Prefix = "goapi:"
)

func GetDefaultKey(name string) string{
	return Prefix + name
}

func GetSqlKey(sql string) string{
	w := md5.New()
	io.WriteString(w, sql)
	md5sql := fmt.Sprintf("%x", w.Sum(nil))
	return Prefix + "sql:" + md5sql
}
