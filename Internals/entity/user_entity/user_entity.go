package user_entity

import (
	"context"
	"renebizelli/go-leilao/Internals/internal_error"

	"github.com/google/uuid"
)

type User struct {
	ID   string
	Name string
}

type UserRepositoryInterface interface {
	FindUserByID(ctx context.Context, userId string) (*User, *internal_error.InternalError)
	CreateUser(ctx context.Context, user *User) *internal_error.InternalError
}

func CreateUser(name string) *User {
	user := &User{
		ID:   uuid.NewString(),
		Name: name,
	}

	return user
}
