package queries

import (
	"context"

	"github.com/edanko/users-api/pkg/decorator"
	"github.com/edanko/users-api/pkg/logs"
)

type ListUsersRequest struct {
	Group *string
}

type ListUsersHandler decorator.QueryHandler[ListUsersRequest, []User]

type ListUsersReadModel interface {
	List(
		ctx context.Context,
		group *string,
	) ([]User, error)
}

type listUsersHandler struct {
	readModel ListUsersReadModel
}

func (h listUsersHandler) Handle(
	ctx context.Context,
	query ListUsersRequest,
) ([]User, error) {
	return h.readModel.List(
		ctx,
		query.Group,
	)
}

func NewListUsersHandler(
	readModel ListUsersReadModel,
	logger logs.Logger,
	metricsClient decorator.MetricsClient,
) ListUsersHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators[ListUsersRequest, []User](
		listUsersHandler{readModel: readModel},
		// logger,
		metricsClient,
	)
}
