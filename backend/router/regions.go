package router

import (
	"evelp/model"

	"github.com/gin-gonic/gin"
)

func regions(c *gin.Context) {
	factions, err := model.GetRegions()
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.JSON(200, factions)
}
