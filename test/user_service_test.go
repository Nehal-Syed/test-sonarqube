package test

import (
	"testing"

	"test-sonarqube/model"
	"test-sonarqube/repository"
	"test-sonarqube/service"
)

func setupTestService() service.UserService {
	repo := repository.NewInMemoryUserRepository()
	return service.NewUserService(repo)
}

func TestUserService_CreateUser(t *testing.T) {
	svc := setupTestService()

	tests := []struct {
		name     string
		userName string
		email    string
		age      int
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid user",
			userName: "John Doe",
			email:    "john@example.com",
			age:      30,
			wantErr:  false,
		},
		{
			name:     "invalid name - too short",
			userName: "J",
			email:    "jane@example.com",
			age:      25,
			wantErr:  true,
			errMsg:   "name must be between 2 and 100 characters",
		},
		{
			name:     "invalid email format",
			userName: "Jane Doe",
			email:    "invalid-email",
			age:      25,
			wantErr:  true,
			errMsg:   "email must be a valid format",
		},
		{
			name:     "invalid age - too young",
			userName: "Bob Smith",
			email:    "bob@example.com",
			age:      16,
			wantErr:  true,
			errMsg:   "age must be between 18 and 120",
		},
		{
			name:     "invalid age - too old",
			userName: "Alice Johnson",
			email:    "alice@example.com",
			age:      150,
			wantErr:  true,
			errMsg:   "age must be between 18 and 120",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := svc.CreateUser(tt.userName, tt.email, tt.age)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && user == nil {
				t.Error("CreateUser() returned nil for valid user")
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("CreateUser() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	svc := setupTestService()

	// Create a test user
	created, _ := svc.CreateUser("Test User", "test@example.com", 25)

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "existing user",
			id:      created.ID,
			wantErr: false,
		},
		{
			name:    "non-existing user",
			id:      999,
			wantErr: true,
		},
		{
			name:    "invalid id",
			id:      0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := svc.GetUser(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && user == nil {
				t.Error("GetUser() returned nil for existing user")
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	svc := setupTestService()

	// Create a test user
	created, _ := svc.CreateUser("Original Name", "original@example.com", 25)

	newName := "Updated Name"
	newEmail := "updated@example.com"
	newAge := 30

	tests := []struct {
		name    string
		id      int
		updates model.UpdateUser
		wantErr bool
	}{
		{
			name: "update name only",
			id:   created.ID,
			updates: model.UpdateUser{
				Name: &newName,
			},
			wantErr: false,
		},
		{
			name: "update email only",
			id:   created.ID,
			updates: model.UpdateUser{
				Email: &newEmail,
			},
			wantErr: false,
		},
		{
			name: "update age only",
			id:   created.ID,
			updates: model.UpdateUser{
				Age: &newAge,
			},
			wantErr: false,
		},
		{
			name:    "update non-existing user",
			id:      999,
			updates: model.UpdateUser{},
			wantErr: true,
		},
		{
			name: "update with invalid age",
			id:   created.ID,
			updates: model.UpdateUser{
				Age: func() *int { i := 150; return &i }(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated, err := svc.UpdateUser(tt.id, tt.updates)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && updated == nil {
				t.Error("UpdateUser() returned nil for valid update")
			}
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	svc := setupTestService()

	// Create a test user
	created, _ := svc.CreateUser("Delete User", "delete@example.com", 25)

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "existing user",
			id:      created.ID,
			wantErr: false,
		},
		{
			name:    "already deleted user",
			id:      created.ID,
			wantErr: true,
		},
		{
			name:    "non-existing user",
			id:      999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.DeleteUser(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
