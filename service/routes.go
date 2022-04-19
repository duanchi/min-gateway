package service

import (
	"fmt"
	"github.com/duanchi/min-gateway/mapper"
	"github.com/duanchi/min-gateway/storage"
	"github.com/duanchi/min-gateway/types/request"
	"github.com/duanchi/min-gateway/types/response"
	"github.com/duanchi/min/abstract"
	util2 "github.com/duanchi/min/util"
	"strings"
)

type Route struct {
	abstract.Bean

	StorageKey          string                       `value:"GATEWAY:ROUTES"`
	KEY                 string                       `value:"routes"`
	RouteStorage        *storage.RouteStorage        `bean:"autowired"`
	RouteRewriteStorage *storage.RouteRewriteStorage `bean:"autowired"`
}

func (this *Route) Init() {}

func (this *Route) GetAll() []response.RouteResponse {
	rawArray := this.RouteStorage.GetAllSorted()
	routeArray := []response.RouteResponse{}
	rewriteMap := this.RouteRewriteStorage.GetAllGroupByRouteId()

	for _, value := range rawArray {
		route := response.RouteResponse{
			Id: value.RouteId,
			Url: response.RouteUrl{
				Type:  mapper.CONSTANT.URL_TYPE[value.UrlType],
				Match: value.Pattern,
			},
			Method:           strings.Split(value.Methods, ","),
			ServiceId:        value.ServiceId,
			Authorize:        mapper.CONSTANT.IS_AUTHORIZE[value.NeedAuthorize],
			AuthorizePrefix:  value.AuthorizePrefix,
			AuthorizeTypeKey: value.AuthorizeTypeKey,
			CustomToken:      mapper.CONSTANT.IS_CUSTOM_TOKEN[value.CustomToken],
			Rewrite:          map[string]string{},
			Order:            value.Sort,
		}

		if value, has := rewriteMap[value.RouteId]; has {
			for _, rewrite := range value {
				route.Rewrite[rewrite.Pattern] = rewrite.Rewrite
			}
		}

		routeArray = append(routeArray, route)
	}
	return routeArray
}

func (this *Route) Add(route request.RouteRequest) {

	route.Id = util2.GenerateUUID().String()
	sort := this.RouteStorage.GetLastSort()
	fmt.Println("SORT", sort)
	route.Order = sort

	this.RouteStorage.Add(mapper.Route{
		RouteId:          route.Id,
		Pattern:          route.Url.Match,
		UrlType:          mapper.CONSTANT.URL_TYPE_REVERSE[route.Url.Type],
		Methods:          strings.Join(route.Method, ","),
		ServiceId:        route.ServiceId,
		NeedAuthorize:    mapper.CONSTANT.IS_AUTHORIZE_REVERSE[route.Authorize],
		AuthorizePrefix:  route.AuthorizePrefix,
		AuthorizeTypeKey: route.AuthorizeTypeKey,
		CustomToken:      mapper.CONSTANT.IS_CUSTOM_TOKEN_REVERSE[route.CustomToken],
		Description:      route.Description,
		Sort:             route.Order,
	})

	if len(route.Rewrite) > 0 {
		rewrites := []mapper.RouteRewrite{}
		for pattern, rewrite := range route.Rewrite {
			rewrites = append(rewrites, mapper.RouteRewrite{
				RouteId: route.Id,
				Pattern: pattern,
				Rewrite: rewrite,
			})
		}
		this.RouteRewriteStorage.AddList(rewrites)
	}
}

func (this *Route) Modify(id string, route request.RouteRequest) {

	rawRoute, ok := this.RouteStorage.GetFromDB(id)

	if ok {
		updated := this.RouteStorage.Update(mapper.Route{
			Id:               rawRoute.Id,
			RouteId:          id,
			Pattern:          route.Url.Match,
			UrlType:          mapper.CONSTANT.URL_TYPE_REVERSE[route.Url.Type],
			Methods:          strings.Join(route.Method, ","),
			ServiceId:        route.ServiceId,
			NeedAuthorize:    mapper.CONSTANT.IS_AUTHORIZE_REVERSE[route.Authorize],
			AuthorizePrefix:  route.AuthorizePrefix,
			AuthorizeTypeKey: route.AuthorizeTypeKey,
			CustomToken:      mapper.CONSTANT.IS_CUSTOM_TOKEN_REVERSE[route.CustomToken],
			Description:      route.Description,
			Sort:             route.Order,
		})

		if updated {
			if len(route.Rewrite) > 0 {
				rewrites := []mapper.RouteRewrite{}
				for pattern, rewrite := range route.Rewrite {
					rewrites = append(rewrites, mapper.RouteRewrite{
						RouteId: id,
						Pattern: pattern,
						Rewrite: rewrite,
					})
				}
				this.RouteRewriteStorage.RemoveByRouteId(id)
				this.RouteRewriteStorage.AddList(rewrites)
			}
		}
	}
}

func (this *Route) Sort(orders []string) {
	this.RouteStorage.Sort(orders)
}

func (this *Route) Remove(id string) {
	this.RouteStorage.Remove(id)
	this.RouteRewriteStorage.RemoveByRouteId(id)
}

func (this *Route) Reload() {
	this.RouteStorage.DataToCache()
	this.RouteRewriteStorage.DataToCache()
}
