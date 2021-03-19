package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type Redis struct {
	name string
}

//func (db *Redis) ModelName() string {
//	return db.name
//}

func New(modelName string) *Redis {
	return &Redis{name: modelName}
}

//func (db *Redis) Do(cmd string, args ...interface{}) Result {
//	return db.Client().Do(cmd, args...)
//}

func (db *Redis) Client() Client {
	return getGlobalClient()
}

func (db *Redis) handleModelKey(key string) string {
	return fmt.Sprintf("%s_%s", db.name, key)
}

func (db *Redis) Get(key string) Result {
	return db.Client().Do("GET", db.handleModelKey(key))
}

func (db *Redis) Set(key string, val interface{}) Result {
	return db.Client().Do("SET", db.handleModelKey(key), val)
}

func (db *Redis) SetWithExpire(key string, val interface{}, seconds interface{}) Result {
	return db.Client().Do("SETEX", db.handleModelKey(key), seconds, val)
}

func (db *Redis) Incr(key string) Result {
	return db.Client().Do("INCR", db.handleModelKey(key))
}

func (db *Redis) Smembers(key string) Result {
	return db.Client().Do("SMEMBERS", db.handleModelKey(key))
}

func (db *Redis) Hset(key string, hkey string, value interface{}) Result {
	return db.Client().Do("HSET", db.handleModelKey(key), hkey, value)
}

func (db *Redis) HGet(key string, hkey string) Result {
	return db.Client().Do("HGET", db.handleModelKey(key), hkey)
}

func (db *Redis) Delete(key string) Result {
	return db.Client().Do("DEL", db.handleModelKey(key))
}

func (db *Redis) Exists(key string) Result {
	return db.Client().Do("EXISTS", db.handleModelKey(key))
}

func (db *Redis) Expire(key string, seconds interface{}) Result {
	return db.Client().Do("EXPIRE", db.handleModelKey(key), seconds)
}

func (db *Redis) Pipe(method func(c Conn) error) (result Result) {
	conn := Conn{db.Client().GetConn()}
	defer conn.self.Close()
	err := method(conn)
	if err != nil {
		return Result{nil, err}
	}
	err = conn.self.Flush()
	if err != nil {
		return Result{nil, err}
	}
	result.reply, result.err = conn.self.Receive()
	return result
}

type Conn struct {
	self redis.Conn
}

func (conn Conn) Send(command string, args ...interface{}) error {
	return conn.self.Send(command, args...)
}
