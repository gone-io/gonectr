package user

import (
	"test/internal/interface/entity"
	"test/internal/pkg/utils"

	"github.com/gone-io/gone"
)

type iUser struct {
	gone.Flag
}

const key = "user"

func (s *iUser) CreateUser(user *entity.User) error {
	utils.Put(key, user)
	return nil
}

func (s *iUser) GetUserList() ([]*entity.User, error) {
	list := utils.Get(key)

	out := make([]*entity.User, 0)
	for _, v := range list {
		out = append(out, v.(*entity.User))
	}

	return out, nil
}
