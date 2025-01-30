package auth

import (
	"Lessons/internal/user"
	"testing"
)

type MockUserRepository struct {
}

func (repo *MockUserRepository) Create(u *user.User) (*user.User, error) {
	return &user.User{
		Email: initialEmail,
	}, nil
}

func (repo *MockUserRepository) FindByEmail(email string) (user *user.User, err error) {
	return nil, nil
}

const initialEmail = "test@test.com"

func TestRegisterSuccess(t *testing.T) {

	authService := NewUserService(&MockUserRepository{})
	email, err := authService.Register(initialEmail, "0", "test")
	if err != nil {
		t.Fatal(err)
	}
	if email != initialEmail {
		t.Fatalf("%s does not match with  %s", email, initialEmail)
	}
}
