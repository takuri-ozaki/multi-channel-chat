package internal

type User struct {
	name string
}

func NewUser(name string) User {
	return User{name: name}
}