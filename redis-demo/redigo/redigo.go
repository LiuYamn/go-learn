package redigo

import (
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Redigo struct {
	pool *redis.Pool
}

func NewRedisPool(host, password string, dbNum int, maxIdle, maxActive int, idleTimeout time.Duration) Redigo {
	redisConn := &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
		Dial: func() (conn redis.Conn, e error) {
			return newRedisPoolFunction(host, password, dbNum)
		},
		TestOnBorrow: testOnBorrowFunction,
	}
	redigo := Redigo{
		pool: redisConn,
	}
	return redigo
}

func (r *Redigo) TestConn() error {
	err := r.Set("test_redis_pool_init", 1, 1)
	if err != nil {
		return err
	}
	return nil
}

func testOnBorrowFunction(c redis.Conn, t time.Time) error {
	_, err := c.Do("PING")
	return err
}

func newRedisPoolFunction(host, password string, dbNum int) (redis.Conn, error) {
	c, err := redis.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	if password != "" {
		if _, err := c.Do("AUTH", password); err != nil {
			err = c.Close()
			if err != nil {
				return nil, err
			}
		}
	}
	_, err = c.Do("SELECT", dbNum)
	if err != nil {
		var _ = c.Close()
		return nil, err
	}
	return c, err
}

func (r *Redigo) Set(key string, data interface{}, time int) error {
	conn := r.pool.Get()
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}
	if time == -1 {
		_, err = conn.Do("EXPIRE", key, 999)
		if err != nil {
			return err
		}
		_, err = conn.Do("PERSIST", key)
		if err != nil {
			return err
		}
	} else {
		_, err = conn.Do("EXPIRE", key, time)
		if err != nil {
			return err
		}
	}
	var _ = conn.Close()
	return nil
}

func (r *Redigo) Exist(key string) bool {
	conn := r.pool.Get()
	exist, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	var _ = conn.Close()
	return exist
}

func (r *Redigo) Get(key string) ([]byte, error) {
	conn := r.pool.Get()
	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}
	var _ = conn.Close()
	return reply, nil
}

func (r *Redigo) Delete(key string) (bool, error) {
	conn := r.pool.Get()
	reply, err := redis.Bool(conn.Do("DEL", key))
	if err != nil {
		return false, err
	}
	var _ = conn.Close()
	return reply, nil
}

func (r *Redigo) DeletesLike(key string) error {
	conn := r.pool.Get()
	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}
	for _, key := range keys {
		_, err = r.Delete(key)
		if err != nil {
			return err
		}
	}
	var _ = conn.Close()
	return nil
}
func (r *Redigo) AddSetMember(setName string, member string) (bool, error) {
	conn := r.pool.Get()
	reply, err := redis.Bool(conn.Do("SADD", setName, member))
	if err != nil {
		return false, err
	}
	var _ = conn.Close()
	return reply, nil
}

func (r *Redigo) ExistSetMember(setName string, member string) (bool, error) {
	conn := r.pool.Get()
	reply, err := redis.Bool(conn.Do("SISMEMBER", setName, member))
	if err != nil {
		return false, err
	}
	var _ = conn.Close()
	return reply, nil
}

func (r *Redigo) DeleteSetMember(setName string, member string) (bool, error) {
	conn := r.pool.Get()
	reply, err := redis.Bool(conn.Do("SREM", setName, member))
	if err != nil {
		return false, err
	}
	var _ = conn.Close()
	return reply, nil
}

func (r *Redigo) GetSetMembers(setName string) ([]string, error) {
	conn := r.pool.Get()
	reply, err := redis.Strings(conn.Do("SMEMBERS", setName))
	if err != nil {
		return []string{}, err
	}
	var _ = conn.Close()
	return reply, nil
}
