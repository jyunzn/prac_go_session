package session

import (
	"errors"
	"sync"

	uuid "github.com/satori/go.uuid"
)

type MemorySessionMgr struct {
	mgr    map[string] *MemorySession
	rwlock sync.RWMutex
}

func NewMemorySessionMgr() *MemorySessionMgr {
	return &MemorySessionMgr{
		mgr: make(map[string] *MemorySession, 1024),
	}
}

func (m *MemorySessionMgr) Init(options ...string) error {
	return nil
}

func (m *MemorySessionMgr) CreateSession() (sessionId string, err error) {

	m.rwlock.Lock()
	defer m.rwlock.Unlock()

	u1 := uuid.NewV4()
	sessionId = u1.String()

	ms := NewMemorySession(sessionId)
	m.mgr[sessionId] = ms

	return sessionId, nil
}

func (m *MemorySessionMgr) Get(sessionId string) (Session, error) {
	m.rwlock.Lock()
	defer m.rwlock.Unlock()

	ms, ok := m.mgr[sessionId]
	if !ok {
		return ms, errors.New("查無此 id")
	}
	return ms, nil
}
