package initial

import (
	"evelp/config/global"

	"github.com/panjf2000/ants/v2"
)

func initAnts() error {
	p, err := ants.NewPool(3)
	if err != nil {
		return err
	}

	global.ANTS = p
	return nil
}
