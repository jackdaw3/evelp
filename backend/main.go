package main

import (
	"evelp/config/global"
	"evelp/initial"
	"evelp/log"
	"evelp/router"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	pprof.Register(engine)
	router.LoadRouter(engine)

	engine.Run(global.Conf.App.ServerPort)
}

func init() {
	if err := initial.Init(); err != nil {
		log.Fatal(err)
	}
}
