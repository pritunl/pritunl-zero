package mhandlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/audit"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/utils"
)

type auditsData struct {
	Audits []*audit.Audit `json:"audits"`
	Count  int64          `json:"count"`
}

func auditsGet(c *gin.Context) {
	if demo.IsDemo() {
		data := &auditsData{
			Audits: demo.Audits,
			Count:  int64(len(demo.Audits)),
		}

		c.JSON(200, data)
		return
	}

	db := c.MustGet("db").(*database.Database)

	page, _ := strconv.ParseInt(c.Query("page"), 10, 0)
	pageCount, _ := strconv.ParseInt(c.Query("page_count"), 10, 0)

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
