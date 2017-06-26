package handlers

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/user"
	"github.com/pritunl/pritunl-zero/utils"
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

func userGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	userId, ok := utils.ParseObjectId(c.Param("user_id"))
	if !ok {
		c.AbortWithStatus(400)
		return
	}

	usr, err := user.Get(db, userId)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, usr)
}

func userPut(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	data := &userData{}

	userId, ok := utils.ParseObjectId(c.Param("user_id"))
	if !ok {
		c.AbortWithStatus(400)
		return
	}

	err := c.Bind(data)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	usr, err := user.Get(db, userId)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	usr.Username = data.Username
	usr.Roles = data.Roles
	usr.Administrator = data.Administrator
	usr.Permissions = data.Permissions

	fields := set.NewSet(
		"username",
		"roles",
		"administrator",
		"permissions",
	)

	if data.Password != "" {
		err = usr.SetPassword(data.Password)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		fields.Add("password")
	}

	errData, err := usr.Validate(db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = usr.CommitFields(db, fields)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, usr)
}

func usersGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	pageStr := c.Query("page")
	page, _ := strconv.Atoi(pageStr)
	pageCountStr := c.Query("page_count")
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
