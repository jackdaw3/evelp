package initial

import (
	"evelp/config/global"
)

func Init() error {
	if err := config(); err != nil {
		return err
	}

	if err := database(); err != nil {
		return err
	}

	if global.Conf.MySQL.AutoMigrate {
		if err := autoMigrate(); err != nil {
			return err
		}
	}

	if err := initRedis(); err != nil {
		return err
	}

	if err := initAnts(); err != nil {
		return err
	}

	if global.Conf.Data.RefreshStaticData {
		if err := staticData(); err != nil {
			return err
		}
	}

	if global.Conf.Data.RefreshRemoteData {
		if err := remoteData(); err != nil {
			return err
		}
	}

	return nil
}
