package service

import (
	"errors"

	"test-sonarqube/model"
	"test-sonarqube/repository"
)

// UserService defines the interface for user business logic
type UserService interface {
	CreateUser(name, email string, age int) (*model.User, error)
	GetUser(id int) (*model.User, error)
	ListUsers() ([]*model.User, error)
	UpdateUser(id int, updates model.UpdateUser) (*model.User, error)
	DeleteUser(id int) error
}

// UserServiceImpl implements UserService
type UserServiceImpl struct {
	repo repository.UserRepository
}

// NewUserService creates a new user service instance
func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{
		repo: repo,
	}
}

// CreateUser creates a new user with validation
func (s *UserServiceImpl) CreateUser(name, email string, age int) (*model.User, error) {
	// Validate input
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	user := &model.User{
		Name:  name,
		Email: email,
		Age:   age,
	}

	// Validate user model
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// Check if email already exists
	if s.repo.ExistsByEmail(email) {
		return nil, errors.New("email already exists")
	}

	// Create user
	err := s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserServiceImpl) GetUser(id int) (*model.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user id")
	}

	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ListUsers returns all users
func (s *UserServiceImpl) ListUsers() ([]*model.User, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUser updates an existing user
func (s *UserServiceImpl) UpdateUser(id int, updates model.UpdateUser) (*model.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user id")
	}

	// Validate updates if provided
	if updates.Name != nil && len(*updates.Name) < 2 {
		return nil, errors.New("name must be at least 2 characters")
	}

	if updates.Age != nil && (*updates.Age < 18 || *updates.Age > 120) {
		return nil, errors.New("age must be between 18 and 120")
	}

	user, err := s.repo.Update(id, updates)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser removes a user
func (s *UserServiceImpl) DeleteUser(id int) error {
	if id <= 0 {
		return errors.New("invalid user id")
	}

	err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
