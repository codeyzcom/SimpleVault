package web

import (
	"SimpleVault/internal/vault"
	"github.com/google/uuid"
	"sync"
	"time"
)

const SessionTTL = 15 * time.Minute

type Session struct {
	Vault     *vault.VaultService
	ExpiresAt time.Time
	UserID    string
}

type SessionManager struct {
	mu       sync.RWMutex
	sessions map[string]*Session

	userIndex map[string]string
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions:  make(map[string]*Session),
		userIndex: make(map[string]string),
	}
}

func (sm *SessionManager) Create(userID string, v *vault.VaultService) string {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if oldSessionID, ok := sm.userIndex[userID]; ok {
		if old, ok := sm.sessions[oldSessionID]; ok {
			old.Vault.Wipe()
		}
		delete(sm.sessions, oldSessionID)
	}

	id := uuid.NewString()
	sm.sessions[id] = &Session{
		Vault:     v,
		ExpiresAt: time.Now().Add(SessionTTL),
		UserID:    userID,
	}
	return id
}

func (sm *SessionManager) Get(id string) (*vault.VaultService, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	s, ok := sm.sessions[id]
	if !ok {
		return nil, false
	}

	if time.Now().After(s.ExpiresAt) {
		s.Vault.Wipe()
		delete(sm.sessions, id)
		delete(sm.userIndex, s.UserID)
		return nil, false
	}

	s.ExpiresAt = time.Now().Add(SessionTTL)

	return s.Vault, true
}

func (sm *SessionManager) Delete(id string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if s, ok := sm.sessions[id]; ok {
		s.Vault.Wipe()
		delete(sm.sessions, id)
		delete(sm.userIndex, s.UserID)
	}
}
