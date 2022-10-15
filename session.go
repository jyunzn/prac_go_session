package session

type Session interface {
	Set(key string, val interface{}) error
	Get(key string) (interface{}, error)
	Del(key string) error
	Save() error
}
