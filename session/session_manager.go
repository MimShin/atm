package session

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

type Session struct {
	UserID    string
	SessionID string
}

type SessionManager struct {
	sessions *sync.Map
}

func NewSessionManager() *SessionManager {
	sm := SessionManager{sessions: &sync.Map{}}
	return &sm
}

func (sm *SessionManager) LoadOrStore(userID, password string) (sessionID string) {
	sessionID = Hash(userID + password + time.Now().String())
	t, _ := sm.sessions.LoadOrStore(userID, sessionID)
	return t.(string)
}

func (sm *SessionManager) Delete(userID string) {
	sm.sessions.Delete(userID)
}

func (sm *SessionManager) Get(userID string) (sessionID string, ok bool) {
	s, ok := sm.sessions.Load(userID)
	return s.(string), ok
}

func (sm *SessionManager) Check(userID, sessionID string) bool {
	s, ok := sm.sessions.Load(userID)
	return ok && s.(string) == sessionID
}

func Hash(str string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
}
