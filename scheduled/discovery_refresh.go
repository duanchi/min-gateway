package scheduled

import (
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/event"
)

type DiscoveryRefreshSchedule struct {
	abstract.Scheduled
}

func (this *DiscoveryRefreshSchedule) Run() {
	event.CommitCondition("DISCOVERY.SERVICE", "CACHED")
}
