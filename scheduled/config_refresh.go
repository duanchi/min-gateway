package scheduled

import (
	"github.com/duanchi/min"
	"github.com/duanchi/min-gateway/service"
	"github.com/duanchi/min/abstract"
	"strconv"
	"time"
)

type ConfigRefreshTask struct {
	abstract.Scheduled
	Interval int64 `value:"${Gateway.ConfigRefreshInterval}"`

	RoutesService   *service.Route   `autowired:"true"`
	ServicesService *service.Service `autowired:"true"`
}

func (this *ConfigRefreshTask) Run() {

	if this.Interval <= 0 {
		return
	}

	// log.Log.Info("Configuration auto refresh enabled!")
	min.Log.Info("Configuration auto refresh enabled!, interval every " + strconv.FormatInt(this.Interval, 10) + " second(s)")

	ticker := time.NewTicker(time.Duration(this.Interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		this.ServicesService.Init()
		this.RoutesService.Init()
	}
}
