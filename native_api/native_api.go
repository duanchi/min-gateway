package native_api

import "reflect"

type NativeApiBeanMap map[string]reflect.Value

var NativeApiBeans = NativeApiBeanMap{}
