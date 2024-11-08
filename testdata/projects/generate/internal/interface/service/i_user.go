package service

import "test/internal/interface/entity"

type IUser interface {
	CreateUser(user *entity.User) error
	GetUserList() ([]*entity.User, error)
}
