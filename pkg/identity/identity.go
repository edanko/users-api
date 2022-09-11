package identity

type Identity struct {
	Email  string   `json:"email"`
	Groups []string `json:"groups"`
	Name   string   `json:"name"`
	Login  string   `json:"login"`
}
