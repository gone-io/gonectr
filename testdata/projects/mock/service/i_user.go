package service

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type IUser interface {
	GetUser(id int64) (*User, error)
	GetUserByUsername(username string) (*User, error)
	CreateUser(user *User) error
	UpdateUser(user *User) error
	DeleteUser(id int64) error
	GetUserPage(query *PageQuery[User]) (*Page[User], error)
}
