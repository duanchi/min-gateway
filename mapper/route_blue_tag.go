package mapper

import "time"

type RouteBlueTag struct {
	Id         int64 `xorm:"pk"`
	RouteId    string
	Tag        string
	ServiceId  string
	CreateTime time.Time `xorm:"created"`
	UpdateTime time.Time `xorm:"updated"`
}
