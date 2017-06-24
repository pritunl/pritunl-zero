package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/user"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

type userData struct {
	Id            bson.ObjectId `json:"id"`
	Type          string        `json:"type"`
	Username      string        `json:"username"`
	Password      string        `json:"password"`
	Roles         []string      `json:"roles"`
	Administrator string        `json:"administrator"`
	Permissions   []string      `json:"permissions"`
}

type usersData struct {
	Users []*user.User `json:"users"`
	Count int          `json:"count"`
}

func usersGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	pageStr := c.Query("page")
	page, _ := strconv.Atoi(pageStr)
	pageCountStr := c.Query("page_ount")
	pageCount, _ := strconv.Atoi(pageCountStr)

	query := &bson.M{}

	users, count, err := user.GetAll(db, query, page, pageCount)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	data := &usersData{
		Users: users,
		Count: count,
	}

	c.JSON(200, data)
}
