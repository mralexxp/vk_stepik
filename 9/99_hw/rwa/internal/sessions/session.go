package sessions

import (
	"fmt"
	"sync"
	"time"
)

// срок истечения сессии (3600 сек = 1 час)
var ExpirationSession int64 = 3600 // секунд

/*
для увеличения производительности можно добавить дополнительное поле, содержащее токены у каждого пользователя
что позволит избежать полного перебора хранилища при использовании DestroyByUsername (напр. при смене пароля).
*/
type SessionManager struct {
	// [PrivateKey]username
	store map[string]*Session
	MU    sync.Mutex
}

type Session struct {
	ID     uint64
	Expire int64
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		store: make(map[string]*Session),
		MU:    sync.Mutex{},
	}
}

func (sm *SessionManager) Create(id uint64) (string, error) {
	if id == 0 {
		return "", fmt.Errorf("id is null")
	}

	public, err := GenerateSession()
	if err != nil || public == "" {
		return "", err
	}

	private := GetPrivateKey(public)

	session := &Session{
		ID:     id,
		Expire: time.Now().Unix() + ExpirationSession,
	}

	sm.MU.Lock()
	defer sm.MU.Unlock()
	sm.store[private] = session

	return public, nil
}

func (sm *SessionManager) Check(public string) (uint64, bool) {
	private := GetPrivateKey(public)

	sm.MU.Lock()
	defer sm.MU.Unlock()

	sess, ok := sm.store[private]
	if !ok {
		return 0, false
	}

	if time.Now().Unix() > sess.Expire {
		delete(sm.store, private)
		return 0, false
	}

	return sess.ID, true
}

func (sm *SessionManager) DestroyByToken(public string) (uint64, error) {
	if public == "" {
		return 0, fmt.Errorf("token is empty")
	}

	private := GetPrivateKey(public)

	sm.MU.Lock()
	defer sm.MU.Unlock()
	sess, ok := sm.store[private]
	if !ok {
		return 0, fmt.Errorf("session '%s' not found", public)
	}

	delete(sm.store, private)

	return sess.ID, nil
}

func (sm *SessionManager) DestroyByID(id uint64) (int, error) {
	if id == 0 {
		return 0, fmt.Errorf("username is empty")
	}

	deleted := 0

	sm.MU.Lock()
	defer sm.MU.Unlock()
	for privateKey, sess := range sm.store {
		if sess.ID == id {
			delete(sm.store, privateKey)
			deleted++
		}
	}

	return deleted, nil
}

// Завернуть в горутину, которая триггерит или вызывает этот метод.
// удаляет из памяти все истекшие сессии
func (sm *SessionManager) ClearExpired() int {
	sm.MU.Lock()
	defer sm.MU.Unlock()

	deleted := 0

	now := time.Now().Unix()

	for privateKey, sess := range sm.store {
		if now > sess.Expire {
			delete(sm.store, privateKey)
			deleted++
		}
	}

	return deleted
}
