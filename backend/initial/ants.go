package initial

import (
	"evelp/config/global"

	"github.com/panjf2000/ants/v2"
)

func initAnts() error {
	pool, err := ants.NewPool(1)
	if err != nil {
		return err
	}
	global.Ants = pool

	return nil
}
