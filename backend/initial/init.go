package initial

import (
	"evelp/config/global"
	"evelp/initial/internal/cachedata"
	"evelp/initial/internal/dbdata"
)

func Init() error {
	if err := config(); err != nil {
		return err
	}

	if err := database(); err != nil {
		return err
	}

	if err := initRedis(); err != nil {
		return err
	}

	if err := initAnts(); err != nil {
		return err
	}

	if global.Conf.Data.Local.Refresh {
		if err := dbdata.LocalData(); err != nil {
			return err
		}
	}

	if global.Conf.Data.Remote.Refresh {
		if err := dbdata.RemoteData(); err != nil {
			return err
		}
	}

	if global.Conf.Redis.Refresh {
		if err := cachedata.CacheData(); err != nil {
			return err
		}
	}

	return nil
}
