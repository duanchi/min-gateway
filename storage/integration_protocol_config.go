package storage

import (
	"github.com/duanchi/min"
	cache2 "github.com/duanchi/min-gateway/cache"
	"github.com/duanchi/min-gateway/mapper"
	"github.com/duanchi/min/abstract"
)

type IntegrationProtocolConfigStorage struct {
	abstract.Service

	CacheService *cache2.CacheService `autowired:"true"`

	CACHE_PREFIX string `value:"HASH:INTEGRATION_PROTOCOL"`
}

func (this *IntegrationProtocolConfigStorage) Get(code string) (config mapper.IntegrationProtocolConfig, ok bool) {

	ok = this.CacheService.Get(this.CACHE_PREFIX, code, &config)
	return
}

func (this *IntegrationProtocolConfigStorage) GetList() (configs mapper.IntegrationProtocolConfig, ok bool) {
	ok = this.CacheService.GetList(this.CACHE_PREFIX, &configs)
	return
}

func (this *IntegrationProtocolConfigStorage) DataToCache() {
	var configs []mapper.IntegrationProtocolConfig
	min.Db.Find(&configs)

	for _, config := range configs {
		this.CacheService.Set(this.CACHE_PREFIX, config.Code, config)
	}
}
