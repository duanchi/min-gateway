package scheduled

import (
	"github.com/duanchi/min-gateway/log"
	"github.com/duanchi/min/abstract"
)

type LogSchedule struct {
	abstract.Scheduled

	LogService *log.LogService `bean:"autowired"`
}

func (this *LogSchedule) Run() {
	this.LogService.Record()
}
