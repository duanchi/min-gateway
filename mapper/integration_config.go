package mapper

import (
	"time"
)

type IntegrationConfig struct {
	Id              int64          `json:"id" xorm:"pk"`
	Uuid            string         `json:"uuid"`
	Alias           string `json:"alias"`
	Url             string         `json:"url"`
	AppId           int64          `json:"app_id"`
	Protocol        string            `json:"protocol"`
	DataType        string  `json:"data_type"`
	RequestMethod   string `json:"request_method"`
	RequestTemplate string `json:"request_template"`
	UpdateTime      time.Time   `json:"update_time" xorm:"updated"`
	CreateTime      time.Time   `json:"create_time" xorm:"created"`
}
