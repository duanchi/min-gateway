package console_api

import (
	"github.com/duanchi/min/types"
	"reflect"
)

type ConsoleApiBeanParser struct {
	types.BeanParser
}

func (parser ConsoleApiBeanParser) Parse(tag reflect.StructTag, bean reflect.Value, definition reflect.Type, beanName string) {
	resource := tag.Get("console_api")

	if resource != "" {
		ConsoleApiBeans[resource] = bean
	}
}
