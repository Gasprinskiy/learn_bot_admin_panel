package chanel_bus

import (
	"sync"
	"time"
)

type Chanel[T any] struct {
	Data  T
	Error error
}

type BusChanelSession[T any] struct {
	Chan   chan Chanel[T]
	Expiry time.Time
}

type BusChanel[T any] struct {
	sessions map[string]BusChanelSession[T]
	mu       sync.Mutex
}

func NewBusChanel[T any]() *BusChanel[T] {
	return &BusChanel[T]{
		sessions: make(map[string]BusChanelSession[T]),
	}
}

func (ac *BusChanel[T]) Read(key string) (BusChanelSession[T], bool) {
	ac.mu.Lock()

	val, exists := ac.sessions[key]

	ac.mu.Unlock()
	return val, exists
}

func (ac *BusChanel[T]) Write(key string, data Chanel[T]) bool {
	ac.mu.Lock()
	session, exists := ac.sessions[key]
	ac.mu.Unlock()

	if !exists || time.Now().After(session.Expiry) {
		ac.CleanUp(key)
		return false
	}

	session.Chan <- data
	return true
}

func (ac *BusChanel[T]) Create(key string, expiryDur time.Duration) {
	ac.mu.Lock()
	ac.sessions[key] = BusChanelSession[T]{
		Chan:   make(chan Chanel[T], 1),
		Expiry: time.Now().Add(expiryDur),
	}
	ac.mu.Unlock()

	go ac.sleepCleanUp(key, expiryDur)
}

func (ac *BusChanel[T]) CleanUp(key string) {
	ac.mu.Lock()

	session, exists := ac.sessions[key]
	if !exists {
		return
	}

	close(session.Chan)
	delete(ac.sessions, key)

	ac.mu.Unlock()
}

func (ac *BusChanel[T]) sleepCleanUp(key string, expiryDur time.Duration) {
	time.Sleep(expiryDur + time.Second)
	ac.CleanUp(key)
}
