package user

import (
	"context"
	"time"
)

type Repository interface {
	CreateBulk(ctx context.Context, users []*User) error
	GetByLogin(ctx context.Context, login string) (*User, error)
	GetLastUpdateTime(ctx context.Context) (time.Time, error)
	// List
}
