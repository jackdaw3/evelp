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

	if global.Conf.Data.RefreshLocalData {
		if err := dbdata.LocalData(); err != nil {
			return err
		}
	}

	if global.Conf.Data.RefreshRemoteData {
		if err := dbdata.RemoteData(); err != nil {
			return err
		}
	}

	if err := cachedata.CacheData(); err != nil {
		return err
	}

	return nil
}
