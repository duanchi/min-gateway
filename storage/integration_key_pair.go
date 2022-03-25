package storage

import (
	cache2 "github.com/duanchi/min-gateway/cache"
	"github.com/duanchi/min-gateway/mapper"
	"github.com/duanchi/min/abstract"
	"strconv"
)

type IntegrationKeyPairStorage struct {
	abstract.Service

	CacheService *cache2.CacheService `autowired:"true"`

	CACHE_PREFIX string `value:"HASH:INTEGRATION_KEY_PAIR"`
}

func (this *IntegrationKeyPairStorage) GetList(integrationId int64) (integrationKeyPairs []mapper.IntegrationKeyPair, ok bool) {
	ok = this.CacheService.GetMatch(this.CACHE_PREFIX, strconv.FormatInt(integrationId, 10)+":*", &integrationKeyPairs)
	return
}

func (this *IntegrationKeyPairStorage) DataToCache() {
	var integrationKeyPairs []mapper.IntegrationKeyPair
	heron.Db.Find(&integrationKeyPairs)

	for _, integrationKeyPair := range integrationKeyPairs {
		this.CacheService.Set(this.CACHE_PREFIX, strconv.FormatInt(integrationKeyPair.IntegrationId, 10)+":"+strconv.FormatInt(integrationKeyPair.Id, 10), integrationKeyPair)
	}
}
