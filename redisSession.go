package session

import (
	"sync"
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
)



type RedisSession struct {
	sessionId string
	pool *redis.Pool
	sessionMap map[string] interface {}
		// set 先將數據存儲到 sessionMap 中
		// 後續執行 save 時，再將這個 map 數據一次存儲到 redis
	flag int
		// 如果數據有變動，就標記一下，這樣 save 就能只操作這些需要被更新到 redis 的數據們
	rwlock sync.RWMutex
}

const (
	SessionFlagNone = iota
	SessionFlagModify
)

func NewRedisSession(sessionId string, pool *redis.Pool) (rs *RedisSession) {
	return &RedisSession {
		sessionId: sessionId,
		pool: pool,
		sessionMap: make(map[string] interface {}, 200),
		flag: SessionFlagNone,
	}
}


func (r *RedisSession)Set(key string, val interface{}) error {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()

	r.sessionMap[key] = val

	// 更新標記
	r.flag = SessionFlagModify
	return nil
}

func (r *RedisSession)Get(key string) (interface{}, error) {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()

	val, ok := r.sessionMap[key]
	if !ok {
		return val, errors.New("查無此 key")
	}
	return val, nil
}

func (r *RedisSession)Del(key string) error {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()

	delete(r.sessionMap, key)
	// 更新標記
	r.flag = SessionFlagModify
	return nil
}

func (r *RedisSession)Save() error {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()

	// 判斷數據是否有修改過
	if r.flag != SessionFlagModify {
		return nil
	}

	// 序列化數據
	data, err := json.Marshal(r.sessionMap)
	if err != nil {
		return err
	}

	// 改數據
	conn := r.pool.Get()
	conn.Do("Set", r.sessionId, string(data))

	// 狀態回復
	r.flag = SessionFlagNone

	return nil
}
