package initial

import (
	"evelp/config/global"

	"github.com/panjf2000/ants/v2"
	"github.com/pkg/errors"
)

func initAnts() error {
	pool, err := ants.NewPool(3)
	if err != nil {
		return errors.Wrap(err, "init ants failed")
	}

	global.ANTS = pool
	return nil
}
