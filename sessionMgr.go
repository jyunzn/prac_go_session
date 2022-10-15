package session


type SessionMgr interface {
	Init(options ...string) error
	CreateSession() (sessionId string, err error)
	Get(sessionId string) (Session, error)
}
