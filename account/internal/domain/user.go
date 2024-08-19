package domain

import (
	"context"
)

type User struct {
	ID    string `json:"id" bson:"_id,omitempty"`
	Email string `json:"email" bson:"email"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetById(ctx context.Context, id string) (*User, error)
}
