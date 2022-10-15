package session

import (
	"errors"
)

func Init(provider string, options ...string) (ssMgr SessionMgr, err error) {
	switch provider {
		case "memory":
			ssMgr = NewMemorySessionMgr()
			return
		case "redis":
			if len(options) == 0 {
				err = errors.New("缺少 addr")
				return
			}
			addr := options[0]
			ssMgr = NewRedisSessionMgr(addr)
			return
		default:
			err = errors.New("只支持 memory || redis")
			return
	}
}
