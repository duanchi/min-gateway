package storage

import (
	"encoding/json"
	"fmt"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/log"
	"github.com/duanchi/min/util"
	"github.com/go-redis/redis/v7"
	"reflect"
	"time"
)

type ValuesService struct {
	abstract.Service

	DataPath string `value:"${Gateway.DataPath}"`
	JwtSignatureKey string `value:"${Authorization.SignatureKey}"`
	JwtExpiresIn int64 `value:"${Authorization.Ttl}"`
	Dsn string `value:"${Authorization.Dsn}"`
	random string
	instance *redis.Client
}

func (this *ValuesService) Init () {
	this.random = util.GenerateUUID().String()
	this.Instance()
}

func (this *ValuesService) Instance () *redis.Client {
	if this.instance == nil {
		options, _ := redis.ParseURL(this.Dsn)

		options.MaxRetries = 3
		options.PoolSize = 8

		this.instance = redis.NewClient(options)
		fmt.Printf("Redis %s connected at DB %d!\r\n", options.Addr, options.DB)
	}
	return this.instance
}


func (this *ValuesService) Set (key string, value interface{}, ttl int64) (err error) {
	valueString,_ := json.Marshal(value)
	_, err = this.Instance().Set(
		key,
		valueString,
		time.Duration(ttl) * time.Second,
		).Result()

	return
}

func (this *ValuesService) HSet (key string, field string, value interface{}, ttl int64) (err error) {
	valueString,_ := json.Marshal(value)
	if ttl == -1 {
		_, err = this.Instance().HSet(
			key,
			field,
			valueString,
		).Result()
	} else {
		_, err = this.Instance().HSet(
			key,
			field,
			valueString,
			time.Duration(ttl) * time.Second,
		).Result()
	}


	return
}

func (this *ValuesService) HGet (key string, field string, value interface{}) (has bool, err error) {
	has = false
	valueString, err := this.Instance().HGet(key, field).Result()

	if err != nil {
		return
	}
	has = true
	err = json.Unmarshal([]byte(valueString), value)
	return
}

func (this *ValuesService) HGetAll (key string, value interface{}) (has bool, err error) {
	has = false
	// *value, err = this.Instance().HGetAll(key).Result()

	valueStringList, err := this.Instance().HGetAll(key).Result()

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

func (this *ValuesService) Get (key string, value interface{}) (has bool, err error) {
	has = false
	valueString, err := this.Instance().Get(key).Result()

	if err != nil {
		return
	}
	has = true
	err = json.Unmarshal([]byte(valueString), value)

	return
}

func (this *ValuesService) Remove (key string) {
	this.Instance().Del(key)
}

func (this *ValuesService) HRemove (key string, field string) {
	this.Instance().HDel(key, field)
}

func (this *ValuesService) Has (key string) (has bool) {
	num, _ := this.Instance().Exists(key).Result()

	return num > 0
}

func (this *ValuesService) HHas (key string, field string) (has bool) {
	has, _ = this.Instance().HExists(key, field).Result()
	return
}