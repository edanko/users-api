package ports

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"

	"github.com/edanko/users-api/internal/app"
	"github.com/edanko/users-api/internal/app/commands"
	queries2 "github.com/edanko/users-api/internal/app/queries"
)

type HTTPServer struct {
	app app.Application
}

func (h HTTPServer) SearchUsers(ctx echo.Context, query string) error {
	appUsers, _ := h.app.Queries.SearchUsers.Handle(ctx.Request().Context(), queries2.SearchUsersRequest{
		Name: query,
	})

	return ctx.JSON(http.StatusOK, appUsers)
}

func NewHTTPServer(application app.Application) HTTPServer {
	return HTTPServer{
		app: application,
	}
}

func (h HTTPServer) GetUser(ctx echo.Context, login string) error {
	appUser, err := h.app.Queries.GetUserByLogin.Handle(ctx.Request().Context(), queries2.GetUserByLoginRequest{
		Login: login,
	})
	if err != nil {
		// httperr.RespondWithSlugError(err, w, r)
		return err
	}

	orderResp := GetUserResponse{
		Login:  appUser.Login,
		Name:   appUser.Name,
		Email:  appUser.Email,
		Groups: appUser.Groups,
	}
	return ctx.JSON(http.StatusOK, orderResp)
}

func (h HTTPServer) ListUsers(ctx echo.Context, params ListUsersParams) error {
	appUsers, err := h.app.Queries.ListUsers.Handle(
		ctx.Request().Context(),
		queries2.ListUsersRequest{
			Group: (*string)(params.Group),
		},
	)
	if err != nil {
		// httperr.RespondWithSlugError(err, w, r)
		return err
	}

	respUsers := lo.Map[queries2.User, User](appUsers, func(e queries2.User, _ int) User {
		return User{
			Login:  e.Login,
			Name:   e.Name,
			Email:  e.Email,
			Groups: e.Groups,
		}
	})

	ordersResp := ListUsersResponse{
		Users: respUsers,
	}
	return ctx.JSON(http.StatusOK, ordersResp)
}

func (h HTTPServer) SyncUsers(ctx echo.Context) error {
	cmd := commands.SyncUsers{}

	err := h.app.CommandBus.Send(ctx.Request().Context(), cmd)
	if err != nil {
		// httperr.RespondWithSlugError(err, w, r)
		return err
	}

	ctx.Response().Status = http.StatusNoContent
	return nil

}
