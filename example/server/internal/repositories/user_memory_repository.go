package repositories

import (
	"fmt"
	"sync"
)

// MemoryUserRepository implementa a interface UserRepository
type MemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*User
}

// NewMemoryUserRepository cria uma nova instância do repositório em memória
func NewMemoryUserRepository() UserRepository {
	return &MemoryUserRepository{
		users: make(map[string]*User),
	}
}

func (r *MemoryUserRepository) Save(user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.ID] = user
	return nil
}

func (r *MemoryUserRepository) GetByID(id string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("user with ID %s not found", id)
	}
	return user, nil
}

func (r *MemoryUserRepository) UpdateByID(id, name string) (*User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("cannot update: user %s not found", id)
	}

	user.Name = name
	r.users[id] = user
	return user, nil
}

func (r *MemoryUserRepository) DeleteByID(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[id]; !ok {
		return fmt.Errorf("cannot delete: user %s not found", id)
	}

	delete(r.users, id)
	return nil
}
