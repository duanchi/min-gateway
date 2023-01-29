package bean

import (
	"github.com/duanchi/min-gateway/cache"
	"github.com/duanchi/min-gateway/console_api"
	"github.com/duanchi/min-gateway/dispatcher"
	"github.com/duanchi/min-gateway/event"
	"github.com/duanchi/min-gateway/log"
	"github.com/duanchi/min-gateway/middleware"
	"github.com/duanchi/min-gateway/native_api/authorize"
	"github.com/duanchi/min-gateway/scheduled"
	"github.com/duanchi/min-gateway/service"
	"github.com/duanchi/min-gateway/storage"
)

type Beans struct {
	AuthorizationService service.AuthorizationService
	TokenService         service.TokenService
	RouteService         service.Route
	ServiceService       service.Service

	ValuesService storage.ValuesService

	ServiceInstanceStorage           storage.ServiceInstanceStorage
	ServiceStorage                   storage.ServiceStorage
	RouteStorage                     storage.RouteStorage
	RouteRewriteStorage              storage.RouteRewriteStorage
	RouteBlueTagStorage              storage.RouteBlueTagStorage
	IntegrationConfigStorage         storage.IntegrationConfigStorage
	IntegrationKeyPairStorage        storage.IntegrationKeyPairStorage
	IntegrationProtocolConfigStorage storage.IntegrationProtocolConfigStorage

	CacheService cache.CacheService
	LogService   log.LogService

	RestfulDispatcher               dispatcher.RestfulDispatcher      `route:"/*url" method:"ALL"`
	ConsoleApiRoutesController      console_api.RoutesController      `console_api:"routes/"`
	ConsoleApiServicesController    console_api.ServicesController    `console_api:"services/"`
	ConsoleApiAuthorizeController   console_api.AuthorizeController   `console_api:"authorize/"`
	ConsoleApiIntegrationController console_api.IntegrationController `console_api:"integration/"`
	ConsoleApiDatasourceController  console_api.DatasourceController  `console_api:"datasource/"`

	NativeApiAuthorizeStatusController authorize.StatusController `native_api:"authorize/status"`

	RouterMiddleware        middleware.RouterMiddleware        `bean:"middleware"`
	AuthorizationMiddleware middleware.AuthorizationMiddleware `bean:"middleware"`
	ConsoleApiMiddleware    middleware.ConsoleApiMiddleware    `bean:"middleware"`
	NativeApiMiddleware     middleware.NativeApiMiddleware     `bean:"middleware"`
	// CustomMiddleware middleware.CustomMiddleware `middleware:"true"`

	CacheSchedule            scheduled.CacheSchedule            `scheduled:"@start"`
	AutoInitServiceSchedule  scheduled.AutoInitServiceSchedule  `scheduled:"@start"`
	DiscoveryRefreshSchedule scheduled.DiscoveryRefreshSchedule `scheduled:"@every 30s"`

	DiscoveryEvent        event.DiscoveryEvent        `event:"DISCOVERY.INIT"`
	DiscoveryServiceEvent event.DiscoveryServiceEvent `event:"DISCOVERY.SERVICE"`
}
