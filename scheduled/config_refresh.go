package scheduled

import (
	"fmt"
	"github.com/duanchi/min-gateway/routes"
	"github.com/duanchi/min/abstract"
	"strconv"
	"time"
)

type ConfigRefreshTask struct {
	abstract.Scheduled
	Interval int64 `value:"${Gateway.ConfigRefreshInterval}"`

	RoutesService   *routes.Routes   `autowired:"true"`
	ServicesService *routes.Services `autowired:"true"`
}

func (this *ConfigRefreshTask) Run() {

	if this.Interval <= 0 {
		return
	}

	// log.Log.Info("Configuration auto refresh enabled!")
	fmt.Println("Configuration auto refresh enabled!, interval every " + strconv.FormatInt(this.Interval, 10) + " second(s)")

	ticker := time.NewTicker(time.Duration(this.Interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		this.ServicesService.Init()
		this.RoutesService.Init()
	}
}
