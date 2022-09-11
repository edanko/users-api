package queries

import (
	"context"

	"github.com/edanko/users-api/internal/domain/user"
	"github.com/edanko/users-api/pkg/decorator"
	"github.com/edanko/users-api/pkg/logs"
)

type GetUserByLoginRequest struct {
	Login string
}

type GetUserByLoginHandler decorator.QueryHandler[GetUserByLoginRequest, User]

type GetUserByLoginReadModel interface {
	GetByLogin(ctx context.Context, login string) (*user.User, error)
}

type getUserByLoginHandler struct {
	readModel GetUserByLoginReadModel
}

func (h getUserByLoginHandler) Handle(
	ctx context.Context,
	query GetUserByLoginRequest,
) (User, error) {
	k, err := h.readModel.GetByLogin(ctx, query.Login)
	if err != nil {
		return User{}, err
	}

	return User{
		Login:  k.Login(),
		Name:   k.Name(),
		Email:  k.Email(),
		Groups: k.Groups(),
	}, nil
}

func NewGetUserByLoginHandler(
	readModel GetUserByLoginReadModel,
	logger logs.Logger,
	metricsClient decorator.MetricsClient,
) GetUserByLoginHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators[GetUserByLoginRequest, User](
		getUserByLoginHandler{readModel: readModel},
		// logger,
		metricsClient,
	)
}
