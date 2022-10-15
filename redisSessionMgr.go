package session

import (
	"sync"
	"time"
	"errors"

	uuid "github.com/satori/go.uuid"
	"github.com/gomodule/redigo/redis"
)

type RedisSessionMgr struct {
	mgr map[string] *RedisSession
	rwlock   sync.RWMutex
	address string  // redis adderss
	// pw      string  // 密碼
	pool    *redis.Pool
}

func NewRedisSessionMgr(addr string) *RedisSessionMgr {
	return &RedisSessionMgr{
		mgr: make(map[string] *RedisSession, 1024),
		address: addr,
	}
}

func connPool(addr string) (*redis.Pool) {
	return &redis.Pool {
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
		MaxIdle: 64,
		MaxActive: 1000,
		IdleTimeout: 100 * time.Second,
	}
}

func (r *RedisSessionMgr)Init(options ...string) (err error) {
	r.pool = connPool(r.address)
	return
}

func (r *RedisSessionMgr)CreateSession() (sessionId string, err error) {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()

	u1 := uuid.NewV4()
	sessionId = u1.String()

	rs := NewRedisSession(sessionId, r.pool)
	r.mgr[sessionId] = rs
	return
}

func (r *RedisSessionMgr)Get(sessionId string) (Session, error) {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()

	rs, ok := r.mgr[sessionId]
	if !ok {
		return rs, errors.New("查無此 seesion")
	}
	return rs, nil
}

