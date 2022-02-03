package initial

import (
	"evelp/config/global"
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func config() error {
	if err := initEnv(); err != nil {
		return err
	}

	if err := initViper(); err != nil {
		return err
	}

	initLogrus()

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

func initLogrus() {
	if global.Conf.App.LogLevel == "info" {
		log.SetLevel(log.InfoLevel)
	}

	if global.Conf.App.LogLevel == "debug" {
		log.SetLevel(log.DebugLevel)
	}
}
