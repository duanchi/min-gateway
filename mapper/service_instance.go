package mapper

import "time"

type ServiceInstance struct {
	Id            int64  `xorm:"pk"`
	InstanceId    string // instance uuid
	GrayFlag      int
	OnlineFlag    int
	Weight        int    // 权重
	Healthy       int    // 是否健康 1健康, 0 不健康
	Url           string // 服务URL
	ServiceId     string
	EphemeralFlag int       // 是否临时实例
	StaticFlag    int       // 是否内置实例
	CreateTime    time.Time `xorm:"created"`
	UpdateTime    time.Time `xorm:"updated"`
	CreateType    int
}
