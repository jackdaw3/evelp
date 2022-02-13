package initial

import (
	"evelp/config/global"
	"evelp/log"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func config() error {
	log.Init()

	if err := initEnv(); err != nil {
		return err
	}

	if err := initViper(); err != nil {
		return err
	}

	log.SetLevel(global.Conf.App.LogLevel)

	return nil
}

func initEnv() error {
	env := os.Getenv("ENV")
	if env == "" {
		log.Warn("no ENV environment variable find on this machine")
		global.ENV = "local"
	} else {
		global.ENV = env
	}
	log.Infof("Env: %s", global.ENV)

	var err error
	global.WORKSPACE, err = os.Getwd()
	if err != nil {
		return errors.Wrap(err, "get workspace failed")
	}

	log.Infof("Workspace: %s", global.WORKSPACE)
	return nil
}

func initViper() error {
	viper.AddConfigPath(global.WORKSPACE + "/config")
	viper.SetConfigName("application-" + global.ENV)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&global.Conf); err != nil {
		return err
	}

	return nil
}
