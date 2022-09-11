package commands

import (
	"context"
	"fmt"

	"github.com/edanko/users-api/internal/adapters"
	"github.com/edanko/users-api/internal/app/events"
	"github.com/edanko/users-api/internal/config"
	"github.com/edanko/users-api/internal/domain/user"
	"github.com/edanko/users-api/pkg/decorator"
	"github.com/edanko/users-api/pkg/errors"
	"github.com/edanko/users-api/pkg/logs"
)

type SyncUsers struct{}

type SyncUsersHandler decorator.CommandHandler

type syncUsersHandler struct {
	config     *config.Config
	eventBus   eventBus
	repository user.Repository
}

func NewSyncUsersHandler(
	config *config.Config,
	eventBus eventBus,
	repo user.Repository,
	logger logs.Logger,
	metricsClient decorator.MetricsClient,
) SyncUsersHandler {
	if repo == nil {
		panic("nil userRepo")
	}

	return decorator.ApplyCommandDecorators(
		syncUsersHandler{
			config:     config,
			eventBus:   eventBus,
			repository: repo,
		},
		// logger,
		metricsClient,
	)
}

func (h syncUsersHandler) HandlerName() string {
	return "sync-users"
}

func (h syncUsersHandler) NewCommand() any {
	return &SyncUsers{}
}

func (h syncUsersHandler) Handle(ctx context.Context, _ any) error {
	ldap := adapters.NewLDAPRepository(h.config)

	users, err := ldap.ListUsers()
	if err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-get-ldap-users")
	}

	err = h.repository.CreateBulk(ctx, users)
	if err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-save-users")
	}
	t, err := h.repository.GetLastUpdateTime(ctx)
	if err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-get-last-update-time")
	}

	fmt.Println(t)
	err = h.eventBus.Publish(ctx, events.UsersSynced{})
	if err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-publish-event")
	}

	return nil
}
