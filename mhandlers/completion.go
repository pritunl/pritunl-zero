package mhandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/aggregate"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/utils"
)

func completionGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	cmpl, err := aggregate.GetCompletion(db, primitive.NilObjectID)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	c.JSON(200, cmpl)
}
