package mapper

import "time"

type SystemTokens struct {
	Id string `xorm:"pk"`
	Expiretime int64
	UserId string
	RefreshId string
	AuthorizeType string
	More map[string]interface{}
	Updatetime time.Time `xorm:"updated"`
	Createtime time.Time `xorm:"created"`
}