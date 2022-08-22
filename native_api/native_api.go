package native_api

import (
	"github.com/duanchi/min/server/types"
)

type NativeApiBeanMap map[string]types.RestfulRoute

var NativeApiBeans = NativeApiBeanMap{}
