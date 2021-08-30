package config

import (
	"github.com/duanchi/min/types"
)

type Config struct {
	types.Config `yaml:",inline"`
	Gateway struct {
		DataPath string `yaml:"data_path" default:"${GATEWAY_DATA_PATH:./data}"`
		CustomMiddleware struct {
			Enabled bool `yaml:"enabled" default:"${GATEWAY_CUSTOM_MIDDLEWARE_ENABLED:false}"`
		} `yaml:"custom_middleware"`
		NativeApi struct {
			Prefix string `yaml:"prefix" default:"${GATEWAY_NATIVE_API_PREFIX:/_api}"`
			AccessToken string `yaml:"access_token" default:"${GATEWAY_NATIVE_API_ACCESS_TOKEN}"`
		} `yaml:"native_api"`
	} `yaml:"gateway"`
}