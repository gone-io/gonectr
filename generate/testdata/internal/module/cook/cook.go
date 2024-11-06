package cook

import (
	"test/internal/interface/entity"
	"test/internal/pkg/utils"

	"github.com/gone-io/gone"
)

type iCook struct {
	gone.Flag
}

const key = "cook"

func (s *iCook) CreateCook(cook *entity.Cook) error {
	utils.Put(key, cook)
	return nil
}
func (s *iCook) GetCookList() ([]*entity.Cook, error) {
	list := utils.Get(key)
	out := make([]*entity.Cook, 0)
	for _, v := range list {
		out = append(out, v.(*entity.Cook))
	}

	return out, nil
}

type IFood struct {
	gone.Flag
}

func (s *IFood) GetMyFood() string {
	return "food"
}
