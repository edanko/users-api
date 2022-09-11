package adapters

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"

	"github.com/edanko/users-api/internal/config"
	domain "github.com/edanko/users-api/internal/domain/user"
)

type LDAPRepository struct {
	config *config.Config
}

func NewLDAPRepository(config *config.Config) *LDAPRepository {
	return &LDAPRepository{
		config: config,
	}
}

func (r *LDAPRepository) ListUsers() ([]*domain.User, error) {
	conn, err := ldap.DialURL(r.config.LDAP.URL)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	err = conn.Bind(r.config.LDAP.User, r.config.LDAP.Password)
	if err != nil {
		return nil, err
	}

	searchRequest := ldap.NewSearchRequest(
		r.config.LDAP.UsersDN+","+r.config.LDAP.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		r.config.LDAP.UsersFilter,
		[]string{"uid", "cn", "mail"},
		nil,
	)

	searchResult, err := conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	users := make([]*domain.User, 0, len(searchResult.Entries))

	for _, u := range searchResult.Entries {
		groupsSearchRequest := ldap.NewSearchRequest(
			r.config.LDAP.GroupsDN+","+r.config.LDAP.BaseDN,
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			fmt.Sprintf(r.config.LDAP.GroupsFilter, ldap.EscapeFilter(u.DN)),
			[]string{"cn"},
			nil,
		)

		groupsSearchResult, err := conn.Search(groupsSearchRequest)
		if err != nil {
			return nil, err
		}

		groups := make([]string, 0, len(groupsSearchResult.Entries))

		for _, g := range groupsSearchResult.Entries {
			if len(g.Attributes) == 0 {
				break
			}

			groups = append(groups, g.Attributes[0].Values...)
		}

		users = append(users, domain.NewUser(
			u.GetAttributeValue("uid"),
			u.GetAttributeValue("cn"),
			u.GetAttributeValue("mail"),
			groups,
		))
	}

	return users, nil
}
