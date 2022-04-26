package mapper

import (
	"time"
)

type Integration struct {
	Id              int64 `xorm:"pk"`
	IntegrationId   string
	Alias           string
	Url             string
	Protocol        string
	DataType        string
	RequestMethod   string
	RequestTemplate string
	UpdateTime      time.Time `xorm:"updated"`
	CreateTime      time.Time `xorm:"created"`
}
