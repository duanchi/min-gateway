package main

import (
	"github.com/duanchi/min-gateway/bean"
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

	Beans bean.Beans
}{
	Config: types.Config{
		BeanParsers: bean.BeanParsers,
	},
}
