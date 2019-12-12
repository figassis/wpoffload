package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"time"

	redis "gopkg.in/redis.v3"
)

func loadRedis() error {
	//opts := redis.Options{Addr: os.Getenv("REDIS"), Password: "", DB: 0}
	opts := redis.Options{Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")), Password: "", DB: 0}
	Log(fmt.Sprintf("Loading Redis: %v", opts))
	cache = redis.NewClient(&opts)
	return Log(cache.Ping().Err())
}

func GetCache(key string, result interface{}) (err error) {
	if key == "" {
		return errors.New("Cache key cannot be empty")
	}
	key = cacheKey + "/" + key
	//Log(fmt.Errorf("KEY: %s", key))
	if reflect.ValueOf(result).Kind() != reflect.Ptr {
		return errors.New("Please provide a pointer")
	}

	obj, err := cache.Get(key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return errors.New("Record not found")
		}
		return
	}

	switch string(obj) {
	case "[]", "{}", "":
		return errors.New("Cache miss")
	}

	go cache.Set(key, obj, cacheTTL)
	if err = Log(json.Unmarshal(obj, result)); err != nil {
		return
	}

	return
}

func DeleteCache(key string) (err error) {
	key = cacheKey + "_" + key
	_, err = cache.Del(key).Result()
	if err == redis.Nil {
		return errors.New("Key not found")
	} else if err != nil {
		return errors.New("Could not remove key")
	}

	return
}

func Cache(key string, obj interface{}) (err error) {
	if key == "" {
		return errors.New("Key cannot be empty")
	}
	key = cacheKey + "_" + key

	if obj == nil {
		return
	}

	//Log(fmt.Sprintf("Caching: %v", obj))
	jsonObj, err := json.Marshal(obj)
	if err != nil {
		return
	}
	return cache.Set(key, jsonObj, cacheTTL).Err()
}

func CacheTTL(key string, obj interface{}, ttl time.Duration) (err error) {
	key = cacheKey + "_" + key

	if obj == nil {
		return
	}

	jsonObj, err := json.Marshal(obj)
	if err != nil {
		return
	}
	return cache.Set(key, jsonObj, ttl).Err()
}
