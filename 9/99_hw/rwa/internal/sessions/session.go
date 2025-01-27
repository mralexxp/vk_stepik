package sessions

import "time"

var ExpirationSession int64 = 3600

type SessionManager struct {
	// [PrivateKey]username
	store map[string]*Session
}

type Session struct {
	Username string
	Expire   int64
}

func NewSessionManager() *SessionManager {
	return &SessionManager{store: make(map[string]*Session)}
}

func (sm *SessionManager) Create(username string) (string, error) {
	public, private, err := GenerateSession(username)
	if err != nil {
		return "", err
	}

	session := &Session{
		Username: username,
		Expire:   time.Now().Unix() + ExpirationSession,
	}

	sm.store[private] = session

	return public, nil
}
