package console_api

import (
	types2 "github.com/duanchi/min/server/types"
	"github.com/duanchi/min/types"
	"github.com/duanchi/min/util"
	"reflect"
)

type ConsoleApiBeanParser struct {
	types.BeanParser
}

func (parser ConsoleApiBeanParser) Parse(tag reflect.StructTag, bean reflect.Value, definition reflect.Type, beanName string) {
	if util.IsBeanKind(tag, "console_api") {

		resource := tag.Get("console_api")
		resourceKey := tag.Get("key")
		if resourceKey == "" {
			resourceKey = "id"
		}
		ConsoleApiBeans[resource] = types2.RestfulRoute{
			Value:       bean,
			ResourceKey: resourceKey,
		}
	}
}
