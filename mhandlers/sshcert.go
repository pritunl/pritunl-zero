package mhandlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/ssh"
	"github.com/pritunl/pritunl-zero/utils"
)

type sshcertsData struct {
	Certificates []*ssh.Certificate `json:"certificates"`
	Count        int64              `json:"count"`
}

func sshcertsGet(c *gin.Context) {
	if demo.IsDemo() {
		data := &sshcertsData{
			Certificates: demo.Sshcerts,
			Count:        int64(len(demo.Sshcerts)),
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

	certs, count, err := ssh.GetCertificates(db, userId, page, pageCount)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	data := &sshcertsData{
		Certificates: certs,
		Count:        count,
	}

	c.JSON(200, data)
}
