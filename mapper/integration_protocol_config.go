package mapper

import "time"

type IntegrationProtocolConfig struct {
	Id			int64		`xorm:"pk"`
	Code		string
	Type		int
	LibraryLink	string
	CreateTime	time.Time	`xorm:"created"`
}