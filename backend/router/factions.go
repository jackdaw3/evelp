package router

import (
	"evelp/model"

	"github.com/gin-gonic/gin"
)

func Factions(c *gin.Context) {
	factions, err := model.GetFactions()
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.JSON(200, factions)
}
