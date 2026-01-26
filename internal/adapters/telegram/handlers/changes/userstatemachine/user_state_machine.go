package userstatemachine

import "sync"

type UserStatesMachine struct {
	userStates map[int64]string
	mu         *sync.Mutex
}

func (u *UserStatesMachine) SetValue(key int64, value string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.userStates[key] = value
}

func (u *UserStatesMachine) Clear(key int64) {
	u.mu.Lock()
	defer u.mu.Unlock()
	delete(u.userStates, key)
}

func (u *UserStatesMachine) GetValue(key int64) (string, bool) {
	u.mu.Lock()
	defer u.mu.Unlock()
	v, ok := u.userStates[key]
	return v, ok
}

func New() *UserStatesMachine {
	return &UserStatesMachine{
		userStates: make(map[int64]string),
		mu:         &sync.Mutex{},
	}
}
