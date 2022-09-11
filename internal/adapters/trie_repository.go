package adapters

import (
	"github.com/edanko/users-api/internal/app/queries"
	"github.com/edanko/users-api/pkg/tr"
)

type TrieRepository struct {
	trie *tr.Trie[queries.User]
}

func NewTrieRepository() *TrieRepository {
	return &TrieRepository{
		trie: tr.New[queries.User](),
	}
}

func (r *TrieRepository) Insert(data []queries.User) {
	for _, user := range data {
		r.trie.Insert(tr.Entry[queries.User]{
			Name:  user.Name,
			Entry: &user,
		})
		r.trie.Insert(tr.Entry[queries.User]{
			Name:  user.Login,
			Entry: &user,
		})
		r.trie.Insert(tr.Entry[queries.User]{
			Name:  user.Email,
			Entry: &user,
		})
	}
}

func (r *TrieRepository) Search(query string) []queries.User {
	return r.trie.SearchAll(query)
}
