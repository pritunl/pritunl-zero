package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/cookie"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/session"
	"github.com/pritunl/pritunl-zero/user"
)

type authData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func authSessionPost(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	data := &authData{}

	err := c.Bind(data)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	usr, err := user.FindUsername(db, user.Local, data.Username)
	if err != nil {
		switch err.(type) {
		case *database.NotFoundError:
			c.JSON(401, &errorData{
				Error:   "auth_invalid",
				Message: "Authencation credentials are invalid",
			})
			break
		default:
			c.AbortWithError(500, err)
		}
		return
	}

	valid := usr.CheckPassword(data.Password)
	if !valid {
		c.JSON(401, &errorData{
			Error:   "auth_invalid",
			Message: "Authencation credentials are invalid",
		})
		return
	}

	cook := cookie.New(c)

	_, err = cook.NewSession(db, usr.Id, true)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.Status(200)
}

func logoutGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	sess := c.MustGet("session").(*session.Session)

	err := sess.Remove(db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.Redirect(302, "/")
}
