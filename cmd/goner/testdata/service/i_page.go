package service

type Page[T any] struct {
	Total int64 `json:"total"`
	List  []*T  `json:"list"`
}

type PageQuery[T any] struct {
	Params   T   `form:"params"`
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

type IPage[T any] interface {
	Get(query PageQuery[T]) (*Page[T], error)
	Create(data T) (T, error)
}
