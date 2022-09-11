package queries

import (
	"context"

	"github.com/edanko/users-api/pkg/decorator"
	"github.com/edanko/users-api/pkg/logs"
)

type SearchUsersRequest struct {
	Name string
}

type SearchUsersHandler decorator.QueryHandler[SearchUsersRequest, []string]

type SearchUsersReadModel interface {
	Search(query string) []string
}

type searchUsersHandler struct {
	readModel SearchUsersReadModel
}

func (h searchUsersHandler) Handle(
	ctx context.Context,
	query SearchUsersRequest,
) ([]string, error) {
	return h.readModel.Search(query.Name), nil
}

func NewSearchUsersHandler(
	readModel SearchUsersReadModel,
	logger logs.Logger,
	metricsClient decorator.MetricsClient,
) SearchUsersHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators[SearchUsersRequest, []string](
		searchUsersHandler{readModel: readModel},
		// logger,
		metricsClient,
	)
}
