package main

import (
	"evelp/config/global"
	"evelp/cron"
	"evelp/initial"
	"evelp/router"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func main() {
	cron.Cron()

	engine := gin.Default()
	router.LoadRouter(engine)
	engine.Run(global.Conf.App.ServerPort)

}

func init() {
	if err := initial.Init(); err != nil {
		log.Fatal(err)
	}
}
