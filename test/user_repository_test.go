package test

import (
	"testing"
	
	"test-sonarqube/model"
	"test-sonarqube/repository"
)

func TestInMemoryUserRepository_Create(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	
	tests := []struct {
		name    string
		user    *model.User
		wantErr bool
	}{
		{
			name: "valid user",
			user: &model.User{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   30,
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			user: &model.User{
				Name:  "Jane Doe",
				Email: "john@example.com",
				Age:   25,
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.user.ID == 0 {
				t.Error("Create() didn't set ID")
			}
		})
	}
}

func TestInMemoryUserRepository_GetByID(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	
	// Create a test user
	user := &model.User{
		Name:  "Test User",
		Email: "test@example.com",
		Age:   25,
	}
	repo.Create(user)
	
	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "existing user",
			id:      user.ID,
			wantErr: false,
		},
		{
			name:    "non-existing user",
			id:      999,
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("GetByID() returned nil for existing user")
			}
		})
	}
}

func TestInMemoryUserRepository_Update(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	
	// Create a test user
	user := &model.User{
		Name:  "Original Name",
		Email: "original@example.com",
		Age:   25,
	}
	repo.Create(user)
	
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
			name: "update all fields",
			id:   user.ID,
			updates: model.UpdateUser{
				Name:  &newName,
				Email: &newEmail,
				Age:   &newAge,
			},
			wantErr: false,
		},
		{
			name:    "update non-existing user",
			id:      999,
			updates: model.UpdateUser{},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated, err := repo.Update(tt.id, tt.updates)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tt.updates.Name != nil && updated.Name != *tt.updates.Name {
					t.Errorf("Update() name = %v, want %v", updated.Name, *tt.updates.Name)
				}
			}
		})
	}
}