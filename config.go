package main

import (
	"github.com/duanchi/min-gateway/dispatcher"
	"github.com/duanchi/min-gateway/middleware"
	"github.com/duanchi/min-gateway/native_api"
	"github.com/duanchi/min-gateway/routes"
	"github.com/duanchi/min-gateway/service"
	"github.com/duanchi/min-gateway/service/storage"
	_interface "github.com/duanchi/min/interface"
	"github.com/duanchi/min/types"
)

var Config = struct {
	types.Config `yaml:",inline"`
	Gateway struct {
		DataPath string `yaml:"data_path" default:"${GATEWAY_DATA_PATH:./data}"`
		CustomMiddleware struct {
			Enabled bool `yaml:"enabled" default:"${GATEWAY_CUSTOM_MIDDLEWARE_ENABLED:false}"`
		} `yaml:"custom_middleware"`
		NativeApi struct {
			Prefix string `yaml:"prefix" default:"${GATEWAY_NATIVE_API_PREFIX:/_api}"`
			AccessToken string `yaml:"access_token"`
		} `yaml:"native_api"`
	} `yaml:"gateway"`

	Authorization struct{
		Ttl int64 `yaml:"ttl" default:"${AUTHORIZATION_TTL:7200}"`
		SignatureKey string `yaml:"signature_key" default:"${AUTHORIZATION_SIGNATURE_KEY}"`
		Dsn string `yaml:"dsn" default:"${AUTHORIZATION_DSN:redis://127.0.0.1:6379/0}"`
		DefaultSingleton bool `yaml:"default_singleton"`
	} `yaml:"authorization"`

	Beans struct {
		AuthorizationService service.AuthorizationService
		TokenService         service.TokenService
		StorageService       storage.StorageService
		ValuesService       storage.ValuesService

		RestfulDispatcher dispatcher.RestfulDispatcher `route:"/*url" method:"ALL"`
		NativeApiRoutesController native_api.RoutesController `native_api:"routes/"`
		NativeApiServicesController native_api.ServicesController `native_api:"services/"`
		NativeApiAuthorizeController native_api.AuthorizeController `native_api:"authorize/"`
		Routes routes.Routes
		Service routes.Services

		RouterMiddleware middleware.RouterMiddleware `middleware:"true"`
		AuthorizationMiddleware middleware.AuthorizationMiddleware `middleware:"true"`
		NativeApiMiddleware middleware.NativeApiMiddleware `middleware:"true"`
		// CustomMiddleware middleware.CustomMiddleware `middleware:"true"`
	}
}{
	Config: types.Config{
		BeanParsers: []_interface.BeanParserInterface{
			&native_api.NativeApiBeanParser{},
		},
	},
}