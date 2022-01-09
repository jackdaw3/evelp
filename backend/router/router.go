package router

import "github.com/gin-gonic/gin"

func LoadRouter(e *gin.Engine) {
	e.GET("/factions", Factions)
	e.GET("/corporations", Corporations)
}
