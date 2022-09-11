package adapters

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	ent2 "github.com/edanko/users-api/internal/adapters/ent"
	user2 "github.com/edanko/users-api/internal/adapters/ent/user"
	"github.com/edanko/users-api/internal/app/queries"
	"github.com/edanko/users-api/internal/domain/user"
)

type UserRepository struct {
	client *ent2.Client
}

var _ user.Repository = (*UserRepository)(nil)

func NewUserRepository(c *ent2.Client) *UserRepository {
	return &UserRepository{
		client: c,
	}
}

func (r *UserRepository) GetLastUpdateTime(ctx context.Context) (time.Time, error) {
	user, err := r.client.User.Query().Order(ent2.Desc("updated_at")).First(ctx)
	if err != nil {
		return time.Time{}, err
	}

	return user.UpdatedAt, nil
}

func (r *UserRepository) CreateBulk(
	ctx context.Context,
	users []*user.User,
) error {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return err
	}

	bulk := make([]*ent2.UserCreate, len(users))
	for i, user := range users {
		bulk[i] = r.client.User.Create().
			SetLogin(user.Login()).
			SetName(user.Name()).
			SetEmail(user.Email()).
			SetGroups(user.Groups())
	}

	err = tx.User.CreateBulk(bulk...).
		OnConflict(sql.ConflictColumns(user2.FieldLogin)).
		UpdateNewValues().
		Exec(ctx)
	if err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("%w: %v", err, rerr)
		}
		return err
	}

	return tx.Commit()
}

func (r *UserRepository) GetByLogin(
	ctx context.Context,
	login string,
) (*user.User, error) {
	e, err := r.client.User.
		Query().
		Where(user2.LoginEQ(login)).
		Only(ctx)
	if ent2.IsNotFound(err) {
		return nil, user.ErrUserNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to get actual user by login")
	}

	return r.unmarshalUser(e), nil
}

func (r *UserRepository) List(
	ctx context.Context,
	group *string,
) ([]queries.User, error) {
	userQuery := r.client.Debug().User.Query().Order(ent2.Asc(user2.FieldLogin))

	if group != nil {
		userQuery = userQuery.Where(func(s *sql.Selector) {
			s.Where(
				sqljson.ValueContains(user2.FieldGroups, *group))
		})
	}

	es, err := userQuery.All(ctx)
	if err != nil {
		return nil, err
	}

	return r.userModelsToQuery(es), nil
}

// func (r *UserRepository) marshalUser(k *domain.User) *ent.User {
// 	return &ent.User{
// 		Login:  k.Login(),
// 		Name:   k.Name(),
// 		Email:  k.Email(),
// 		Groups: k.Groups(),
// 	}
// }

func (r *UserRepository) unmarshalUser(e *ent2.User) *user.User {
	return user.NewUser(
		e.Login,
		e.Name,
		e.Email,
		e.Groups,
	)
}

func (r *UserRepository) userModelsToQuery(es []*ent2.User) []queries.User {
	return lo.Map[*ent2.User, queries.User](es, func(e *ent2.User, _ int) queries.User {
		return queries.User{
			Login:  e.Login,
			Name:   e.Name,
			Email:  e.Email,
			Groups: e.Groups,
		}
	})
}
