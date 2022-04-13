package cache

import (
	"context"
	"encoding/json"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/config"
	"github.com/go-redis/redis/v8"
	"reflect"
	"strings"
	"time"
)

type CacheService struct {
	abstract.Service

	Engine *redis.Client
	ctx    context.Context
	DSN    string `value:"${Cache.Dsn}"`
}

func (this *CacheService) FlushDB() {
	this.Engine.FlushDB(this.ctx)
}

func (this *CacheService) Init() {
	dsn := config.Get("Cache.Dsn").(string)
	options, _ := redis.ParseURL(dsn)
	this.Engine = redis.NewClient(options)
	this.ctx = context.Background()
}

func (this *CacheService) Get(prefix string, key string, value interface{}) (ok bool) {
	valueString, error := this.Engine.HGet(this.ctx, prefix, key).Result()
	if error != nil || valueString == "" {
		ok = false
	} else {
		json.Unmarshal([]byte(valueString), value)
		ok = true
	}
	return
}

func (this *CacheService) GetList(prefix string, value interface{}) (ok bool) {
	result, _, error := this.Engine.HScan(this.ctx, prefix, 0, "*", 65535).Result()

	if error == nil {
		for _, v := range result {
			if strings.Contains(v, "{") {
				elem := reflect.New(reflect.TypeOf(value).Elem().Elem()).Interface()
				json.Unmarshal([]byte(v), &elem)
				reflect.ValueOf(value).Elem().Set(reflect.Append(reflect.ValueOf(value).Elem(), reflect.ValueOf(elem).Elem()))
			}
		}
	}

	return
}

func (this *CacheService) GetMatch(prefix string, match string, value interface{}) (ok bool) {
	result, _, error := this.Engine.HScan(this.ctx, prefix, 0, match, 65535).Result()

	if error == nil {
		for _, v := range result {
			if strings.Contains(v, "{") {
				elem := reflect.New(reflect.TypeOf(value).Elem().Elem()).Interface()
				json.Unmarshal([]byte(v), &elem)
				reflect.ValueOf(value).Elem().Set(reflect.Append(reflect.ValueOf(value).Elem(), reflect.ValueOf(elem).Elem()))
			}
		}
	}

	return
}

func (this *CacheService) GetInKeys(prefix string, key string, value interface{}) (ok bool) {
	valueString, error := this.Engine.Get(this.ctx, prefix+":"+key).Result()
	if error != nil || valueString == "" {
		ok = false
	} else {
		json.Unmarshal([]byte(valueString), value)
		ok = true
	}
	return
}

func (this *CacheService) Set(prefix string, key string, value interface{}) {
	v, _ := json.Marshal(value)
	this.Engine.HSet(this.ctx, prefix, key, v).Result()
}

func (this *CacheService) SetToKeys(prefix string, key string, value interface{}) {
	v, _ := json.Marshal(value)
	setArg := redis.SetArgs{
		Mode:     "",
		TTL:      0,
		ExpireAt: time.Time{},
		Get:      false,
		KeepTTL:  true,
	}
	this.Engine.SetArgs(this.ctx, prefix+":"+key, v, setArg).Result()
}

func (this *CacheService) SetWithTTL(prefix string, key string, value interface{}, ttl int) {
	v, _ := json.Marshal(value)
	this.Engine.Set(this.ctx, prefix+":"+key, v, time.Duration(ttl)*time.Second).Result()
}

func (this *CacheService) Increase(prefix string, key string) {
	this.Engine.Incr(this.ctx, prefix+":"+key)
}

func (this *CacheService) DelPrefix(prefix string) {
	this.Engine.Del(this.ctx, prefix+":*").Result()
	this.Engine.HDel(this.ctx, prefix, "*").Result()
}

func (this *CacheService) Del(prefix string, key string) {
	this.Engine.Del(this.ctx, prefix+":"+key).Result()
	this.Engine.HDel(this.ctx, prefix, key).Result()
}
