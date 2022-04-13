package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/config"
	"github.com/duanchi/min/log"
	"github.com/duanchi/min/util"
	"github.com/go-redis/redis/v8"
	"reflect"
	"strings"
	"time"
)

type ValuesService struct {
	abstract.Service

	DataPath        string `value:"${Gateway.DataPath}"`
	JwtSignatureKey string `value:"${Authorization.SignatureKey}"`
	JwtExpiresIn    int64  `value:"${Authorization.Ttl}"`
	Dsn             string `value:"${Authorization.Dsn}"`
	random          string
	instance        *redis.Client
	context         context.Context
}

func (this *ValuesService) Init() {
	this.context = context.Background()
	this.random = util.GenerateUUID().String()
	this.Instance()
}

func (this *ValuesService) useMemoryCache() bool {
	return this.Dsn == ""
}

func (this *ValuesService) Instance() *redis.Client {
	if this.instance == nil {
		dsn := config.Get("Authorization.Dsn").(string)
		options, _ := redis.ParseURL(dsn)

		options.MaxRetries = 3
		options.PoolSize = 8

		this.instance = redis.NewClient(options)
		fmt.Printf("Redis %s connected at DB %d!\r\n", options.Addr, options.DB)
	}
	return this.instance
}

func (this *ValuesService) Set(key string, value interface{}, ttl int64) (err error) {
	valueString, _ := json.Marshal(value)
	_, err = this.Instance().Set(
		this.context,
		key,
		valueString,
		time.Duration(ttl)*time.Second,
	).Result()

	return
}

func (this *ValuesService) HSet(key string, field string, value interface{}, ttl int64) (err error) {
	valueString, _ := json.Marshal(value)
	if ttl == -1 {
		_, err = this.Instance().HSet(
			this.context,
			key,
			field,
			valueString,
		).Result()
	} else {
		_, err = this.Instance().HSet(
			this.context,
			key,
			field,
			valueString,
			time.Duration(ttl)*time.Second,
		).Result()
	}

	return
}

func (this *ValuesService) HGet(key string, field string, value interface{}) (has bool, err error) {
	has = false
	valueString, err := this.Instance().HGet(this.context, key, field).Result()

	if err != nil {
		return
	}
	has = true
	err = json.Unmarshal([]byte(valueString), value)
	return
}

func (this *ValuesService) HGetAll(key string, value interface{}) (has bool, err error) {
	has = false
	// *value, err = this.Instance().HGetAll(key).Result()

	valueStringList, err := this.Instance().HGetAll(this.context, key).Result()

	if err != nil {
		return
	}
	has = true

	for key, valueString := range valueStringList {
		val := reflect.New(reflect.TypeOf(value).Elem().Elem())

		err = json.Unmarshal([]byte(valueString), val.Interface())
		if err != nil {
			log.Log.Fatal(err)
		}
		reflect.ValueOf(value).Elem().SetMapIndex(reflect.ValueOf(key), val.Elem())
	}
	// err = json.Unmarshal([]byte(valueString), value)
	return
}

func (this *ValuesService) Get(key string, value interface{}) (has bool, err error) {
	has = false
	valueString, err := this.Instance().Get(this.context, key).Result()

	if err != nil {
		return
	}
	has = true
	err = json.Unmarshal([]byte(valueString), value)

	return
}

func (this *ValuesService) Remove(key string) {
	this.Instance().Del(this.context, key)
}

func (this *ValuesService) HRemove(key string, field string) {
	this.Instance().HDel(this.context, key, field)
}

func (this *ValuesService) Has(key string) (has bool) {
	num, _ := this.Instance().Exists(this.context, key).Result()

	return num > 0
}

func (this *ValuesService) HHas(key string, field string) (has bool) {
	has, _ = this.Instance().HExists(this.context, key, field).Result()
	return
}

func (this *ValuesService) GetAll(value interface{}) (has bool, err error) {
	keys, err := this.Instance().Keys(this.context, "*").Result()

	if err == nil {
		for _, key := range keys {
			v, kerr := this.Instance().Get(this.context, key).Result()
			if nil == kerr && strings.Contains(v, "{") {
				elem := reflect.New(reflect.TypeOf(value).Elem().Elem()).Interface()
				json.Unmarshal([]byte(v), &elem)
				reflect.ValueOf(value).Elem().SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(elem).Elem())
			}
		}
	}
	return
}
