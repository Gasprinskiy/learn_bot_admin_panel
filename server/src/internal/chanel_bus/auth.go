package chanel_bus

import (
	"learn_bot_admin_panel/internal/entity/profile"
	"sync"
	"time"
)

type SessionChanel struct {
	User  profile.User
	Error error
}

type Session struct {
	Chan   chan SessionChanel
	Expiry time.Time
}

type AuthChan struct {
	sessions map[string]Session
	mu       sync.Mutex
}

func NewAuthChan() *AuthChan {
	return &AuthChan{
		sessions: make(map[string]Session),
	}
}

func (ac *AuthChan) Read(key string) (Session, bool) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	val, exists := ac.sessions[key]

	return val, exists
}

func (ac *AuthChan) Write(key string, data SessionChanel) bool {
	ac.mu.Lock()
	session, exists := ac.sessions[key]

	defer ac.mu.Unlock()

	if !exists || time.Now().After(session.Expiry) {
		ac.CleanUp(key)
		return false
	}

	session.Chan <- data
	return true
}

func (ac *AuthChan) Create(key string, expiryDur time.Duration) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	ac.sessions[key] = Session{
		Chan:   make(chan SessionChanel, 1),
		Expiry: time.Now().Add(expiryDur),
	}
}

func (ac *AuthChan) CleanUp(key string) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	session, exists := ac.sessions[key]
	if !exists {
		return
	}

	close(session.Chan)
	delete(ac.sessions, key)
}
