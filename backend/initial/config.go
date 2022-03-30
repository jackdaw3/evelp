package initial

import (
	"evelp/config/global"
	"evelp/log"
	"os"

	"github.com/gin-gonic/gin"
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

	initGinMode()
	log.SetLevel(global.Conf.App.LogLevel)

	return nil
}

func initEnv() error {
	env := os.Getenv("ENV")
	if env == "" {
		return errors.New("no ENV environment variable find on this machine")
	} else {
		global.Env = env
	}
	log.Infof("Env: %s", global.Env)

	var err error
	global.Workspace, err = os.Getwd()
	if err != nil {
		return errors.Wrap(err, "get workspace failed")
	}

	log.Infof("Workspace: %s", global.Workspace)
	return nil
}

func initViper() error {
	viper.AddConfigPath(global.Workspace + "/config")
	viper.SetConfigName("application-" + global.Env)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&global.Conf); err != nil {
		return err
	}

	return nil
}

func initGinMode() {
	if global.Env == "local" {
		gin.SetMode(gin.DebugMode)
	}

	if global.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
}
