package mhandlers

import (
	"fmt"
	"github.com/dropbox/godropbox/container/set"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/user"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"strings"
	"time"
)

type userData struct {
	Id            bson.ObjectId `json:"id"`
	Type          string        `json:"type"`
	Username      string        `json:"username"`
	Password      string        `json:"password"`
	Roles         []string      `json:"roles"`
	Administrator string        `json:"administrator"`
	Permissions   []string      `json:"permissions"`
	Disabled      bool          `json:"disabled"`
	ActiveUntil   time.Time     `json:"active_until"`
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

	usr.Type = data.Type
	usr.Username = data.Username
	usr.Roles = data.Roles
	usr.Administrator = data.Administrator
	usr.Permissions = data.Permissions
	usr.Disabled = data.Disabled
	usr.ActiveUntil = data.ActiveUntil

	if usr.Disabled {
		usr.ActiveUntil = time.Time{}
	}

	fields := set.NewSet(
		"type",
		"username",
		"roles",
		"administrator",
		"permissions",
		"disabled",
		"active_until",
	)

	if usr.Type == user.Local && data.Password != "" {
		err = usr.SetPassword(data.Password)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		fields.Add("password")
	} else if usr.Type != user.Local && usr.Password != "" {
		usr.Password = ""
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

	event.PublishDispatch(db, "user.change")

	c.JSON(200, usr)
}

func userPost(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	data := &userData{}

	err := c.Bind(data)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	usr := &user.User{
		Type:          data.Type,
		Username:      data.Username,
		Roles:         data.Roles,
		Administrator: data.Administrator,
		Permissions:   data.Permissions,
		Disabled:      data.Disabled,
		ActiveUntil:   data.ActiveUntil,
	}

	if usr.Disabled {
		usr.ActiveUntil = time.Time{}
	}

	if usr.Type == user.Local && data.Password != "" {
		err = usr.SetPassword(data.Password)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
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

	err = usr.Insert(db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	event.PublishDispatch(db, "user.change")

	c.JSON(200, usr)
}

func usersGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	pageStr := c.Query("page")
	page, _ := strconv.Atoi(pageStr)
	pageCountStr := c.Query("page_count")
	pageCount, _ := strconv.Atoi(pageCountStr)

	query := bson.M{}

	username := strings.TrimSpace(c.Query("username"))
	if username != "" {
		query["username"] = &bson.M{
			"$regex":   fmt.Sprintf(".*%s.*", username),
			"$options": "i",
		}
	}

	role := strings.TrimSpace(c.Query("role"))
	if role != "" {
		query["roles"] = role
	}

	typ := strings.TrimSpace(c.Query("type"))
	if typ != "" {
		query["type"] = typ
	}

	administrator := c.Query("administrator")
	switch administrator {
	case "true":
		query["administrator"] = "super"
		break
	case "false":
		query["administrator"] = ""
		break
	}

	disabled := c.Query("disabled")
	switch disabled {
	case "true":
		query["disabled"] = true
		break
	case "false":
		query["disabled"] = false
		break
	}

	users, count, err := user.GetAll(db, &query, page, pageCount)
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

func usersDelete(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	userIds := []bson.ObjectId{}

	err := c.Bind(&userIds)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	errData, err := user.Remove(db, userIds)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	event.PublishDispatch(db, "user.change")

	c.JSON(200, nil)
}
