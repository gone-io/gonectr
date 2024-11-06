package service

import "test/internal/interface/entity"

type ICook interface {
	CreateCook(cook *entity.Cook) error
	GetCookList() ([]*entity.Cook, error)
}
