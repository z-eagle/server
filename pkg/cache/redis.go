package cache

import (
	"bytes"
	"encoding/gob"
	"github.com/zhouqiaokeji/server/pkg/util"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

// RedisStore Redis Store Driver
type RedisStore struct {
	pool *redis.Pool
}

type item struct {
	Value interface{}
}

func serializer(value interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	storeValue := item{
		Value: value,
	}
	err := enc.Encode(storeValue)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func deserializer(value []byte) (interface{}, error) {
	var res item
	buffer := bytes.NewReader(value)
	dec := gob.NewDecoder(buffer)
	err := dec.Decode(&res)
	if err != nil {
		return nil, err
	}
	return res.Value, nil
}

// NewRedisStore Create Redis Store
func NewRedisStore(size int, network, address, password, database string) *RedisStore {
	return &RedisStore{
		pool: &redis.Pool{
			MaxIdle:     size,
			IdleTimeout: 240 * time.Second,
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
			Dial: func() (redis.Conn, error) {
				db, err := strconv.Atoi(database)
				if err != nil {
					return nil, err
				}

				c, err := redis.Dial(
					network,
					address,
					redis.DialDatabase(db),
					redis.DialPassword(password),
				)
				if err != nil {
					util.Log().Warning("Redis Connect Fail：%s", err)
					return nil, err
				}
				return c, nil
			},
		},
	}
}

// Set Save Store Value
func (store *RedisStore) Set(key string, value interface{}, ttl int) error {
	rc := store.pool.Get()
	defer rc.Close()

	serialized, err := serializer(value)
	if err != nil {
		return err
	}

	if rc.Err() != nil {
		return rc.Err()
	}

	if ttl > 0 {
		_, err = rc.Do("SETEX", key, ttl, serialized)
	} else {
		_, err = rc.Do("SET", key, serialized)
	}

	if err != nil {
		return err
	}
	return nil

}

// Get Value
func (store *RedisStore) Get(key string) (interface{}, bool) {
	rc := store.pool.Get()
	defer rc.Close()
	if rc.Err() != nil {
		return nil, false
	}

	v, err := redis.Bytes(rc.Do("GET", key))
	if err != nil || v == nil {
		return nil, false
	}

	finalValue, err := deserializer(v)
	if err != nil {
		return nil, false
	}

	return finalValue, true

}

// Gets Batch Get Value
func (store *RedisStore) Gets(keys []string, prefix string) (map[string]interface{}, []string) {
	rc := store.pool.Get()
	defer rc.Close()
	if rc.Err() != nil {
		return nil, keys
	}

	var queryKeys = make([]string, len(keys))
	for key, value := range keys {
		queryKeys[key] = prefix + value
	}

	v, err := redis.ByteSlices(rc.Do("MGET", redis.Args{}.AddFlat(queryKeys)...))
	if err != nil {
		return nil, keys
	}

	var res = make(map[string]interface{})
	var missed = make([]string, 0, len(keys))

	for key, value := range v {
		decoded, err := deserializer(value)
		if err != nil || decoded == nil {
			missed = append(missed, keys[key])
		} else {
			res[keys[key]] = decoded
		}
	}
	return res, missed
}

// Sets Batch Set Values
func (store *RedisStore) Sets(values map[string]interface{}, prefix string) error {
	rc := store.pool.Get()
	defer rc.Close()
	if rc.Err() != nil {
		return rc.Err()
	}
	var setValues = make(map[string]interface{})

	for key, value := range values {
		serialized, err := serializer(value)
		if err != nil {
			return err
		}
		setValues[prefix+key] = serialized
	}

	_, err := rc.Do("MSET", redis.Args{}.AddFlat(setValues)...)
	if err != nil {
		return err
	}
	return nil

}

// Delete Batch Remove keys
func (store *RedisStore) Delete(keys []string, prefix string) error {
	rc := store.pool.Get()
	defer rc.Close()
	if rc.Err() != nil {
		return rc.Err()
	}

	// 处理前缀
	for i := 0; i < len(keys); i++ {
		keys[i] = prefix + keys[i]
	}

	_, err := rc.Do("DEL", redis.Args{}.AddFlat(keys)...)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAll Clean
func (store *RedisStore) DeleteAll() error {
	rc := store.pool.Get()
	defer rc.Close()
	if rc.Err() != nil {
		return rc.Err()
	}

	_, err := rc.Do("FLUSHDB")

	return err
}
