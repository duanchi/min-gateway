package mapper

import (
	"time"
)

type IntegrationKeyPair struct {
	Id            int64         `json:"id" xorm:"pk"`
	IntegrationId int64         `json:"integration_id"`
	Key           string        `json:"key"`
	Expression    string        `json:"expression"`
	Position      int `json:"position"`
	UpdateTime    time.Time  `json:"update_time" xorm:"updated"`
	CreateTime    time.Time  `json:"create_time" xorm:"created"`
}
