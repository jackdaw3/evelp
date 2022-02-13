package router

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func LoadRouter(e *gin.Engine) {
	e.GET("/faction", faction)
	e.GET("/corporation", corporation)
	e.GET("/region", region)
	e.GET("/offer", offer)
	e.GET("/history", itemHistory)

	pprof.Register(e)
}
