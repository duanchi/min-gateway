package scheduled

import (
	"github.com/duanchi/min-gateway/cache"
	"github.com/duanchi/min-gateway/storage"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/event"
)

type CacheSchedule struct {
	abstract.Scheduled
	ServiceInstanceStorage           *storage.ServiceInstanceStorage           `bean:"autowired"`
	ServiceStorage                   *storage.ServiceStorage                   `bean:"autowired"`
	RouteStorage                     *storage.RouteStorage                     `bean:"autowired"`
	RouteRewriteStorage              *storage.RouteRewriteStorage              `bean:"autowired"`
	IntegrationConfigStorage         *storage.IntegrationConfigStorage         `bean:"autowired"`
	IntegrationKeyPairStorage        *storage.IntegrationKeyPairStorage        `bean:"autowired"`
	IntegrationProtocolConfigStorage *storage.IntegrationProtocolConfigStorage `bean:"autowired"`

	CacheService *cache.CacheService `bean:"autowired"`
}

func (this *CacheSchedule) Run() {
	this.CacheService.FlushDB()
	this.ServiceStorage.DataToCache()
	this.ServiceInstanceStorage.DataToCache()
	this.RouteStorage.DataToCache()
	this.RouteRewriteStorage.DataToCache()
	this.IntegrationConfigStorage.DataToCache()
	this.IntegrationKeyPairStorage.DataToCache()
	this.IntegrationProtocolConfigStorage.DataToCache()
	event.CommitCondition("DISCOVERY.SERVICE", "CACHED")
}
