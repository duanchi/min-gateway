package mapper

import "time"

type Route struct {
	Id               int64 `xorm:"pk"`
	RouteType        int   // 路由类型 0 服务路由, 1 url路由
	Pattern          string
	UrlType          int // 0 path_match, 1 fn_match, 2 regex
	Methods          string
	ServiceId        int64
	NeedAuthorize    int    // 是否需要授权
	AuthorizePrefix  string // 授权因子
	AuthorizeTypeKey string // 区分单例登录的Key
	CustomToken      int    // 是否自定义授权token
	Description      string
	Sort             int
	CreateTime       time.Time `xorm:"created"`
	UpdateTime       time.Time `xorm:"updated"`
}
