package native_api

import (
	types2 "github.com/duanchi/min/server/types"
	"github.com/duanchi/min/types"
	"github.com/duanchi/min/util"
	"reflect"
)

type NativeApiBeanParser struct {
	types.BeanParser
}

func (parser NativeApiBeanParser) Parse(tag reflect.StructTag, bean reflect.Value, definition reflect.Type, beanName string) {
	if util.IsBeanKind(tag, "native_api") {

		resource := tag.Get("native_api")
		resourceKey := tag.Get("key")
		if resourceKey == "" {
			resourceKey = "id"
		}
		NativeApiBeans[resource] = types2.RestfulRoute{
			Value:       bean,
			ResourceKey: resourceKey,
		}
	}
}
