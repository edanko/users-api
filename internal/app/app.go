package app

import (
	"context"

	queries2 "github.com/edanko/users-api/internal/app/queries"
)

type mediator interface {
	Send(ctx context.Context, cmd any) error
}

type Application struct {
	CommandBus mediator
	Queries    Queries
}

type Queries struct {
	ListUsers      queries2.ListUsersHandler
	GetUserByLogin queries2.GetUserByLoginHandler
	SearchUsers    queries2.SearchUsersHandler
}
