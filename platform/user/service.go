package user

import (
	"context"
)

type ListUsersOption struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`

	FilterUserID []string `json:"filter_user_id"`
	FilterUDID   []string `json:"filter_udid"`
}

type Service interface {
	ApplyUser(ctx context.Context, u User) (*User, error)
	ListUsers(ctx context.Context, opt ListUsersOption) ([]User, error)
}

type Store interface {
	User(context.Context, string) (*User, error)
	Save(*User) error
	List() ([]User, error)
}

type UserService struct {
	store Store
}

func New(store Store) *UserService {
	return &UserService{store: store}
}
