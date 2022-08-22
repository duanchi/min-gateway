package storage

import (
	"github.com/duanchi/min"
	cache2 "github.com/duanchi/min-gateway/cache"
	"github.com/duanchi/min-gateway/mapper"
	"github.com/duanchi/min/abstract"
	"strconv"
)

type RouteBlueTagStorage struct {
	abstract.Service

	CacheService *cache2.CacheService `autowired:"true"`
	CACHE_PREFIX string               `value:"HASH:ROUTE_BLUE_TAG"`
}

func (this *RouteBlueTagStorage) GetByRouteId(id string) (tags []mapper.RouteBlueTag) {
	allTags := []mapper.RouteBlueTag{}
	this.CacheService.GetList(this.CACHE_PREFIX, &allTags)

	for _, tag := range allTags {
		if tag.RouteId == id {
			tags = append(tags, tag)
		}
	}

	return
}

func (this *RouteBlueTagStorage) Get(id int64) (tag mapper.RouteBlueTag) {
	this.CacheService.Get(this.CACHE_PREFIX, strconv.FormatInt(id, 10), &tag)
	return
}

func (this *RouteBlueTagStorage) GetAllGroupByRouteId() (tags map[string][]mapper.RouteBlueTag) {
	allTags := []mapper.RouteBlueTag{}
	this.CacheService.GetList(this.CACHE_PREFIX, &allTags)
	tags = map[string][]mapper.RouteBlueTag{}

	for _, tag := range allTags {
		if _, has := tags[tag.RouteId]; !has {
			tags[tag.RouteId] = []mapper.RouteBlueTag{}
		}
		tags[tag.RouteId] = append(tags[tag.RouteId], tag)
	}

	return
}

func (this *RouteBlueTagStorage) Add(tag mapper.RouteBlueTag) (id int64, ok bool) {
	_, err := min.Db.Insert(&tag)
	if err != nil {
		return 0, false
	} else {
		this.DataToCache()
		return tag.Id, true
	}
}

func (this *RouteBlueTagStorage) AddList(tags []mapper.RouteBlueTag) {
	_, err := min.Db.Insert(&tags)
	if err == nil {
		this.DataToCache()
	}
}

func (this *RouteBlueTagStorage) RemoveByRouteId(id string) {
	tag := mapper.RouteBlueTag{}
	min.Db.Where("route_id = ?", id).Delete(tag)
	this.DataToCache()
}

func (this *RouteBlueTagStorage) Remove(id int64) {
	tag := this.Get(id)
	min.Db.Where("code = ?", id).Delete(tag)
	this.DataToCache()
}

func (this *RouteBlueTagStorage) DataToCache() {
	var tags []mapper.RouteBlueTag
	min.Db.Find(&tags)

	this.CacheService.DelPrefix(this.CACHE_PREFIX)

	for _, tag := range tags {
		this.CacheService.Set(this.CACHE_PREFIX, strconv.FormatInt(tag.Id, 10), tag)
	}
}
