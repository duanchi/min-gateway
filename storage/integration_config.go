package storage

import (
	"github.com/duanchi/min"
	cache2 "github.com/duanchi/min-gateway/cache"
	"github.com/duanchi/min-gateway/mapper"
	"github.com/duanchi/min/abstract"
	"text/template"
)

type IntegrationConfigStorage struct {
	abstract.Service

	CacheService *cache2.CacheService `autowired:"true"`

	requestTemplate map[string]*template.Template

	CACHE_PREFIX       string `value:"HASH:INTEGRATION_CONFIG"`
	CACHE_ALIAS_PREFIX string `value:"HASH:INTEGRATION_ALIAS_CONFIG"`
}

func (this *IntegrationConfigStorage) Init() {
	this.requestTemplate = map[string]*template.Template{}
}

func (this *IntegrationConfigStorage) Get(uuid string) (integrationConfig mapper.IntegrationConfig, ok bool) {
	ok = this.CacheService.Get(this.CACHE_PREFIX, uuid, &integrationConfig)
	return
}

/*func (this *IntegrationConfigStorage) GetUuidByAlias (alias string) (uuid string, ok bool) {
	ok = this.CacheService.Get(this.CACHE_ALIAS_PREFIX, alias, &uuid)
	return
}*/

func (this *IntegrationConfigStorage) GetByAlias(alias string) (integrationConfig mapper.IntegrationConfig, ok bool) {
	ok = this.CacheService.Get(this.CACHE_ALIAS_PREFIX, alias, &integrationConfig)

	return
}

func (this *IntegrationConfigStorage) GetTemplate(uuid string) (tmpl *template.Template, ok bool) {
	tmpl, ok = this.requestTemplate[uuid]

	return
}

func (this *IntegrationConfigStorage) DataToCache() {
	var integrationConfigs []mapper.IntegrationConfig
	min.Db.Find(&integrationConfigs)

	for _, integrationConfig := range integrationConfigs {

		templateInstance, err := template.New(integrationConfig.Uuid).Parse(integrationConfig.RequestTemplate)
		if err == nil {
			this.requestTemplate[integrationConfig.Uuid] = templateInstance
		}

		this.CacheService.Set(this.CACHE_PREFIX, integrationConfig.Uuid, integrationConfig)

		if integrationConfig.Alias != "" {
			this.CacheService.Set(this.CACHE_ALIAS_PREFIX, integrationConfig.Alias, integrationConfig)
		}
	}
}
