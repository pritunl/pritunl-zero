package mhandlers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/log"
	"github.com/pritunl/pritunl-zero/utils"
)

type logsData struct {
	Logs  []*log.Entry `json:"logs"`
	Count int64        `json:"count"`
}

func logGet(c *gin.Context) {
	if demo.IsDemo() {
		c.JSON(200, demo.Logs[1])
		return
	}

	db := c.MustGet("db").(*database.Database)

	logId, ok := utils.ParseObjectId(c.Param("log_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	usr, err := log.Get(db, logId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	c.JSON(200, usr)
}

func logsGet(c *gin.Context) {
	if demo.IsDemo() {
		data := &logsData{
			Logs:  demo.Logs,
			Count: int64(len(demo.Logs)),
		}

		c.JSON(200, data)
		return
	}

	db := c.MustGet("db").(*database.Database)

	pageStr := c.Query("page")
	page, _ := strconv.ParseInt(pageStr, 10, 0)
	pageCountStr := c.Query("page_count")
	pageCount, _ := strconv.ParseInt(pageCountStr, 10, 0)

	query := bson.M{}

	message := strings.TrimSpace(c.Query("message"))
	if message != "" {
		query["message"] = &bson.M{
			"$regex":   fmt.Sprintf(".*%s.*", regexp.QuoteMeta(message)),
			"$options": "i",
		}
	}

	level := strings.TrimSpace(c.Query("level"))
	if level != "" {
		query["level"] = level
	}

	logs, count, err := log.GetAll(db, &query, page, pageCount)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	data := &logsData{
		Logs:  logs,
		Count: count,
	}

	c.JSON(200, data)
}
