package mapper

import (
	"time"
)

type IntegrationKeyPair struct {
	Id            int64 ` xorm:"pk"`
	IntegrationId string
	Key           string
	Expression    string
	Position      int
	UpdateTime    time.Time `xorm:"updated"`
	CreateTime    time.Time `xorm:"created"`
}
