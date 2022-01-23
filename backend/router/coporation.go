package router

import (
	"evelp/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func corporation(c *gin.Context) {
	factionId, err := strconv.Atoi(c.Param("factionId"))
	if err != nil {
		c.AbortWithError(500, err)
	}

	corporations, err := model.GetCorporationsByFaction(factionId)
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.JSON(200, corporations)
}
