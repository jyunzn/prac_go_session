package session

import (
	"errors"
	"sync"
)

type MemorySession struct {
	SessionId string
	data      map[string]interface{}
	rwlock    sync.RWMutex
}

func NewMemorySession(id string) *MemorySession {
	return &MemorySession{
		SessionId: id,
		data:      make(map[string]interface{}, 100),
	}
}

func (m *MemorySession) Set(key string, val interface{}) error {

	m.rwlock.Lock()
	defer m.rwlock.Unlock()

	m.data[key] = val

	return nil
}

func (m *MemorySession) Get(key string) (interface{}, error) {
	m.rwlock.Lock()
	defer m.rwlock.Unlock()

	data, ok := m.data[key]
	if !ok {
		return data, errors.New("查無此數據")
	}
	return data, nil
}

func (m *MemorySession) Del(key string) error {

	m.rwlock.Lock()
	defer m.rwlock.Unlock()
	delete(m.data, key)

	return nil
}

func (m *MemorySession) Save() error {
	return nil
}
