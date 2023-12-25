package userrepository

import (
	"context"
	"github.com/google/uuid"
)

//go:generate mockgen -destination mock_user_repository.go -package userrepository . IUserRepository

type IUserRepository interface {
	Create(ctx context.Context, u *User) error
	Get(ctx context.Context, id uuid.UUID) (*User, error)
	GetByUserAndPassword(ctx context.Context, login string) (*User, error)
	List(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByLoginCheck(ctx context.Context, login string) (bool, error)
	GetByEmailCheck(ctx context.Context, email string) (bool, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByLogin(ctx context.Context, login string) (*User, error)
}
