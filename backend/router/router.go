package router

import "github.com/gin-gonic/gin"

func LoadRouter(e *gin.Engine) {
	e.GET("/factions", factions)
	e.GET("/corporations", corporations)
	e.GET("/regions", regions)
}
