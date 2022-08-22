package event

import (
	"github.com/duanchi/min/abstract"
	event2 "github.com/duanchi/min/event"
	"github.com/duanchi/min/types"
)

type DiscoveryEvent struct {
	abstract.Event
}

func (this *DiscoveryEvent) Emit(event types.Event, arguments ...interface{}) {
	event2.CommitCondition("DISCOVERY.SERVICE", "DISCOVERED")
}
