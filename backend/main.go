package main

import (
	"evelp/config/global"
	"evelp/initial"
	"evelp/router"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	engine := gin.Default()
	router.LoadRouter(engine)
	engine.Run(global.Conf.App.ServerPort)
}

func init() {
	if err := initial.Init(); err != nil {
		log.Fatal(err)
	}
}
