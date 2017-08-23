package mhandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/audit"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/utils"
	"strconv"
)

type auditsData struct {
	Audits []*audit.Audit `json:"audits"`
	Count  int            `json:"count"`
}

func auditsGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	page, _ := strconv.Atoi(c.Query("page"))
	pageCount, _ := strconv.Atoi(c.Query("page_count"))

	userId, ok := utils.ParseObjectId(c.Param("user_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	audits, count, err := audit.GetAll(db, userId, page, pageCount)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	data := &auditsData{
		Audits: audits,
		Count:  count,
	}

	c.JSON(200, data)
}
