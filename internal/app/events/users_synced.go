package events

import (
	"context"

	"github.com/rs/zerolog/log"
)

type UsersSynced struct{}

type UsersSyncedHandler struct{}

func (k UsersSyncedHandler) HandlerName() string {
	return "users-synced"
}

func (k UsersSyncedHandler) NewEvent() any {
	return &UsersSynced{}
}

func (k UsersSyncedHandler) Handle(_ context.Context, _ any) error {
	// ctx, cancel := context.WithTimeout(ctx, time.Second*120)
	// defer cancel()

	// e := events.(*UsersSynced)
	// _ = e

	log.Info().Msg("UsersSynced events received")

	return nil
}
