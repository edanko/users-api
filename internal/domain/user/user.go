package user

type User struct {
	login  string
	name   string
	email  string
	groups []string
}

func (k *User) Login() string {
	return k.login
}

func (k *User) Name() string {
	return k.name
}

func (k *User) Email() string {
	return k.email
}

func (k *User) Groups() []string {
	return k.groups
}

func NewUser(
	login string,
	name string,
	email string,
	groups []string,
) *User {
	return &User{
		login:  login,
		name:   name,
		email:  email,
		groups: groups,
	}
}
