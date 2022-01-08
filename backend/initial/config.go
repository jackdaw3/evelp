package initial

import (
	"evelp/configs/global"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func config() error {
	if err := initEnv(); err != nil {
		return err
	}

	if err := initConfig(); err != nil {
		return err
	}

	return nil
}

func initEnv() error {
	env := os.Getenv("ENV")
	if env == "" {
		log.Warn("No ENV environment variable find on this machine!")
		global.ENV = "local"
	} else {
		global.ENV = env
	}
	log.Info("Env: ", global.ENV)

	var err error
	global.WORKSPACE, err = os.Getwd()
	if err != nil {
		return err
	}

	log.Info("Workspace: ", global.WORKSPACE)
	return nil
}

func initConfig() error {

	viper.AddConfigPath(global.WORKSPACE + "/configs")
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
