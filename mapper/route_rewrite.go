package mapper

import "time"

type RouteRewrite struct {
	Id         int64 `xorm:"pk"`
	RouteId    string
	Pattern    string
	Rewrite    string
	CreateTime time.Time `xorm:"created"`
	UpdateTime time.Time `xorm:"updated"`
}
