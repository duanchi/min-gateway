package main

import (
	"github.com/duanchi/min-gateway/cache"
	"github.com/duanchi/min-gateway/console_api"
	"github.com/duanchi/min-gateway/dispatcher"
	"github.com/duanchi/min-gateway/middleware"
	"github.com/duanchi/min-gateway/native_api"
	"github.com/duanchi/min-gateway/native_api/authorize"
	"github.com/duanchi/min-gateway/scheduled"
	"github.com/duanchi/min-gateway/service"
	storage2 "github.com/duanchi/min-gateway/service/storage"
	"github.com/duanchi/min-gateway/storage"
	_interface "github.com/duanchi/min/interface"
	"github.com/duanchi/min/types"
)

var Config = struct {
	types.Config `yaml:",inline"`
	Gateway      struct {
		DataPath              string `yaml:"data_path" default:"${GATEWAY_DATA_PATH:./data}"`
		ConfigRefreshInterval int64  `yaml:"config_refresh_interval" default:"${CONFIG_REFRESH_INTERVAL:0}"`
		ConsoleApi            struct {
			Prefix      string `yaml:"prefix" default:"${GATEWAY_CONSOLE_API_PREFIX:/_api}"`
			AccessToken string `yaml:"access_token"`
		} `yaml:"console_api"`
		NativeApi struct {
			Prefix string `yaml:"prefix" default:"${GATEWAY_NATIVE_API_PREFIX:/native}"`
		} `yaml:"native_api"`
		EnableGrayInstance   bool  `yaml:"enable_gray_instance" default:"${ENABLE_GRAY_INSTANCE:true}"`
		GlobalRequestTimeout int64 `yaml:"global_request_timeout" default:"${GLOBAL_REQUEST_TIMEOUT:30}"`
	} `yaml:"gateway"`

	Authorization struct {
		Ttl              int64  `yaml:"ttl" default:"${AUTHORIZATION_TTL:7200}"`
		SignatureKey     string `yaml:"signature_key" default:"${AUTHORIZATION_SIGNATURE_KEY}"`
		Dsn              string `yaml:"dsn" default:"${AUTHORIZATION_DSN:}"`
		DefaultSingleton bool   `yaml:"default_singleton"`
	} `yaml:"authorization"`

	Beans struct {
		AuthorizationService service.AuthorizationService
		TokenService         service.TokenService
		RouteService         service.Route
		ServiceService       service.Service

		ValuesService storage2.ValuesService

		ServiceInstanceStorage           storage.ServiceInstanceStorage
		ServiceStorage                   storage.ServiceStorage
		RouteStorage                     storage.RouteStorage
		RouteRewriteStorage              storage.RouteRewriteStorage
		IntegrationConfigStorage         storage.IntegrationConfigStorage
		IntegrationKeyPairStorage        storage.IntegrationKeyPairStorage
		IntegrationProtocolConfigStorage storage.IntegrationProtocolConfigStorage

		CacheService cache.CacheService

		RestfulDispatcher               dispatcher.RestfulDispatcher      `route:"/*url" method:"ALL"`
		ConsoleApiRoutesController      console_api.RoutesController      `console_api:"routes/"`
		ConsoleApiServicesController    console_api.ServicesController    `console_api:"services/"`
		ConsoleApiAuthorizeController   console_api.AuthorizeController   `console_api:"authorize/"`
		ConsoleApiIntegrationController console_api.IntegrationController `console_api:"integration/"`
		ConsoleApiDatasourceController  console_api.DatasourceController  `console_api:"datasource/"`

		NativeApiAuthorizeStatusController authorize.StatusController `native_api:"authorize/status"`

		RouterMiddleware        middleware.RouterMiddleware        `middleware:"true"`
		AuthorizationMiddleware middleware.AuthorizationMiddleware `middleware:"true"`
		ConsoleApiMiddleware    middleware.ConsoleApiMiddleware    `middleware:"true"`
		NativeApiMiddleware     middleware.NativeApiMiddleware     `middleware:"true"`
		// CustomMiddleware middleware.CustomMiddleware `middleware:"true"`

		CacheSchedule scheduled.CacheSchedule `scheduled:"@start"`
	}
}{
	Config: types.Config{
		BeanParsers: []_interface.BeanParserInterface{
			&console_api.ConsoleApiBeanParser{},
			&native_api.NativeApiBeanParser{},
		},
	},
}
