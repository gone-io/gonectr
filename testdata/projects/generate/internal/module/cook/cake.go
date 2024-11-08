package cook

import "github.com/gone-io/gone"

type iCake struct {
	gone.Flag
}

func (s *iCake) Eat() error {
	return nil
}
