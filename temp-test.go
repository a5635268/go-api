package main

import (
	"fmt"
	"reflect"
)

func main2() {
	type Time struct {
		Id int `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
		Name string `json:"name"`
		Create_time int `json:"create_time" gorm:"autoCreateTime"`      // 使用秒级时间戳填充创建时间
		Update_time int `json:"update_time" gorm:"autoUpdateTime"`
	}
	time := Time{
		Id : 7,
		Name : "lilei",
		Create_time : 1602658039,
	}

	//reflect.ValueOf(time).Type().Kind()
	m := reflect.ValueOf(time).Kind()

	fmt.Println(m)
}
