package mapper

import "time"

type Service struct {
	Id               int64 `xorm:"pk"`
	Code             string
	Name             string
	LoadBalanceType  int
	NeedHealthCheck  int // 是否需要健康检查 1,需要,0 不需要
	NeedBreak        int // 是否需要熔断, 1需要, 0不需要
	BreakMaxTries    int // 在${fuse_window}时间范围内达到${fuse_max_tries}后触发熔断
	BreakWindow      int
	BreakRestartTime int
	CreateTime       time.Time `xorm:"created"`
	UpdateTime       time.Time `xorm:"updated"`
}
