package bean

import (
	"github.com/duanchi/min-gateway/console_api"
	"github.com/duanchi/min-gateway/native_api"
	_interface "github.com/duanchi/min/interface"
)

var BeanParsers = []_interface.BeanParserInterface{
	&console_api.ConsoleApiBeanParser{},
	&native_api.NativeApiBeanParser{},
}
