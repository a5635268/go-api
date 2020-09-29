package tools

import (
	"fmt"
	"reflect"
)

func Print(e ...interface{}) {
	for _, v := range e {
		fmt.Printf("\n")
		if v == nil {
			fmt.Printf("%T: %#v\n",v,v)
			return
		}
		kind := reflect.TypeOf(v).Kind()
		switch kind {
		case reflect.Struct:
			res,_ :=  StructToJsonStr(v)
			fmt.Printf("struct: %v\n",res)
		case reflect.Map:
			res,_ :=  MapToJsonStr(v)
			fmt.Printf("map: %v\n",res)
		case reflect.Slice:
			res,_ := SliceToJsonStr(v)
			fmt.Printf("slice: %v\n",res)
		default:
			fmt.Printf("%T: %#v\n",v,v)
		}
	}
}
