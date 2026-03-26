package repository

import (
	"errors"
	"sync"
	"time"

	"test-sonarqube/model"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *model.User) error
	GetByID(id int) (*model.User, error)
	GetAll() ([]*model.User, error)
	Update(id int, updates model.UpdateUser) (*model.User, error)
	Delete(id int) error
	ExistsByEmail(email string) bool
}

// InMemoryUserRepository implements UserRepository with in-memory storage
type InMemoryUserRepository struct {
	users  map[int]*model.User
	mu     sync.RWMutex
	nextID int
}

// NewInMemoryUserRepository creates a new instance of in-memory repository
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:  make(map[int]*model.User),
		nextID: 1,
	}
}

// Create adds a new user to the repository
func (r *InMemoryUserRepository) Create(user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validate email uniqueness
	if r.existsByEmailUnsafe(user.Email) {
		return errors.New("email already exists")
	}

	user.ID = r.nextID
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	
	r.users[user.ID] = user
	r.nextID++
	
	return nil
}

// GetByID retrieves a user by ID
func (r *InMemoryUserRepository) GetByID(id int) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	
	return user, nil
}

// GetAll returns all users
func (r *InMemoryUserRepository) GetAll() ([]*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*model.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	
	return users, nil
}

// Update updates an existing user
func (r *InMemoryUserRepository) Update(id int, updates model.UpdateUser) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	// Apply updates
	if updates.Name != nil {
		user.Name = *updates.Name
	}
	if updates.Email != nil {
		// Check email uniqueness if changed
		if *updates.Email != user.Email && r.existsByEmailUnsafe(*updates.Email) {
			return nil, errors.New("email already exists")
		}
		user.Email = *updates.Email
	}
	if updates.Age != nil {
		user.Age = *updates.Age
	}
	
	user.UpdatedAt = time.Now()
	return user, nil
}

// Delete removes a user from the repository
func (r *InMemoryUserRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return errors.New("user not found")
	}
	
	delete(r.users, id)
	return nil
}

// ExistsByEmail checks if a user with given email exists
func (r *InMemoryUserRepository) ExistsByEmail(email string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	return r.existsByEmailUnsafe(email)
}

// existsByEmailUnsafe is an internal method that doesn't acquire locks
func (r *InMemoryUserRepository) existsByEmailUnsafe(email string) bool {
	for _, user := range r.users {
		if user.Email == email {
			return true
		}
	}
	return false
}