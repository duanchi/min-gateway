package native_api

import (
	"github.com/duanchi/min/types"
	"reflect"
)

type NativeApiBeanParser struct {
	types.BeanParser
}

func (parser NativeApiBeanParser) Parse (tag reflect.StructTag, bean reflect.Value, definition reflect.Type, beanName string) {
	resource := tag.Get("native_api")

	if resource != "" {
		NativeApiBeans[resource] = bean
	}
}