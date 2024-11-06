package foods

type Food struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

type IFood interface {
	Create(food *Food) (*Food, error)
	Get(id int) (*Food, error)
	Update(food *Food) (*Food, error)
	Delete(id int) error
}

type IFoodUser interface {
	Create() (*Food, error)
}
