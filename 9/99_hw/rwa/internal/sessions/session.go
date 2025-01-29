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
	Username string
	Expire   int64
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		store: make(map[string]*Session),
		MU:    sync.Mutex{},
	}
}

// создает новую сессию в хранилище
func (sm *SessionManager) Create(username string) (string, error) {
	if username == "" {
		return "", fmt.Errorf("username is empty")
	}

	public, err := GenerateSession(username)
	if err != nil || public == "" {
		return "", err
	}

	private := GetPrivateKey(public)

	session := &Session{
		Username: username,
		Expire:   time.Now().Unix() + ExpirationSession,
	}

	sm.MU.Lock()
	defer sm.MU.Unlock()
	sm.store[private] = session

	return public, nil
}

// проверяет валидность токена и возвращает ok и username, которому принадлежит токен
func (sm *SessionManager) Check(public string) (string, bool) {
	private := GetPrivateKey(public)

	sm.MU.Lock()
	defer sm.MU.Unlock()

	sess, ok := sm.store[private]
	if !ok {
		return "", false
	}

	if time.Now().Unix() > sess.Expire {
		delete(sm.store, private)
		return "", false
	}

	return sess.Username, true
}

// удаляет сессию по публичному токену
func (sm *SessionManager) DestroyByToken(public string) (string, error) {
	if public == "" {
		return "", fmt.Errorf("token is empty")
	}

	private := GetPrivateKey(public)

	sm.MU.Lock()
	defer sm.MU.Unlock()
	sess, ok := sm.store[private]
	if !ok {
		return "", fmt.Errorf("%s not found", public)
	}

	delete(sm.store, private)

	return sess.Username, nil
}

// удаляет сессии пользователя по нику
func (sm *SessionManager) DestroyByUsername(username string) (int, error) {
	if username == "" {
		return 0, fmt.Errorf("username is empty")
	}

	deleted := 0

	sm.MU.Lock()
	defer sm.MU.Unlock()
	for privateKey, sess := range sm.store {
		if sess.Username == username {
			delete(sm.store, privateKey)
			deleted++
		}
	}

	return deleted, nil
}

// TODO: Завернуть в горутину, которая триггерит или вызывает этот метод.
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
