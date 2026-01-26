package timer

import (
	"sync"

	"github.com/google/uuid"
)

type Store struct {
	userTimers map[int64]map[string]struct{}
	mu         *sync.Mutex
}

func NewStore() *Store {
	return &Store{
		userTimers: make(map[int64]map[string]struct{}),
		mu:         &sync.Mutex{},
	}
}

func (t *Store) NewTimer(chatID int64) string {
	t.mu.Lock()
	defer t.mu.Unlock()
	newTimer, _ := uuid.NewUUID()
	if _, ok := t.userTimers[chatID]; !ok {
		t.userTimers[chatID] = make(map[string]struct{})
	}
	t.userTimers[chatID][newTimer.String()] = struct{}{}
	return newTimer.String()
}

func (t *Store) HasTimer(chatID int64, timerID string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if v, ok := t.userTimers[chatID]; ok {
		if _, ok2 := v[timerID]; ok2 {
			return true
		}
	}
	return false
}

func (t *Store) StopTimer(chatID int64, timerID string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if v, ok := t.userTimers[chatID]; ok {
		if _, ok2 := v[timerID]; ok2 {
			delete(t.userTimers[chatID], timerID)
			return true
		}
	}
	return false
}
