package router

import (
	"evelp/model"

	"github.com/gin-gonic/gin"
)

func faction(c *gin.Context) {
	factions, err := model.GetFactions()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, factions)
}
