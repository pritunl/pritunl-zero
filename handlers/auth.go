package handlers

import (
	"crypto/hmac"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"github.com/Sirupsen/logrus"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/auth"
	"github.com/pritunl/pritunl-zero/cookie"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/session"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/user"
	"gopkg.in/mgo.v2/bson"
	"net/url"
	"strings"
)

type authStateData struct {
	Providers []*authStateProviderData `json:"providers"`
}

type authStateProviderData struct {
	Id    bson.ObjectId `json:"id"`
	Type  string        `json:"type"`
	Label string        `json:"label"`
}

func authStateGet(c *gin.Context) {
	data := &authStateData{
		Providers: []*authStateProviderData{},
	}

	for _, provider := range settings.Auth.Providers {
		providerData := &authStateProviderData{
			Id:    provider.Id,
			Type:  provider.Type,
			Label: provider.Label,
		}
		data.Providers = append(data.Providers, providerData)
	}

	c.JSON(200, data)
}

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

	usr, err := user.GetUsername(db, user.Local, data.Username)
	if err != nil {
		switch err.(type) {
		case *database.NotFoundError:
			c.JSON(401, &errortypes.ErrorData{
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
		c.JSON(401, &errortypes.ErrorData{
			Error:   "auth_invalid",
			Message: "Authencation credentials are invalid",
		})
		return
	}

	if usr.Administrator != "super" {
		c.JSON(401, &errortypes.ErrorData{
			Error:   "unauthorized",
			Message: "Not authorized",
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

	if sess != nil {
		err := sess.Remove(db)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
	}

	c.Redirect(302, "/login")
}

func authRequestGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	providerId := bson.ObjectIdHex(c.Query("id"))

	var provider *settings.Provider
	for _, prvidr := range settings.Auth.Providers {
		if prvidr.Id == providerId {
			provider = prvidr
			break
		}
	}

	if provider == nil {
		c.AbortWithStatus(404)
		return
	}

	loc := location.Get(c).String()

	switch provider.Type {
	case auth.Google:
		redirect, err := auth.GoogleRequest(db, loc, provider)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.Redirect(302, redirect)
		return
	case auth.OneLogin, auth.Okta:
		body, err := auth.SamlRequest(db, loc, provider)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.Data(200, "text/html;charset=utf-8", body)
		return
	}

	c.AbortWithStatus(404)
}

func authCallbackGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	query := strings.Split(c.Request.URL.RawQuery, "&sig=")[0]

	params, err := url.ParseQuery(query)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	state := params.Get("state")
	sig := c.Query("sig")

	tokn, err := auth.Get(db, state)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	hashFunc := hmac.New(sha512.New, []byte(tokn.Secret))
	hashFunc.Write([]byte(query))
	rawSignature := hashFunc.Sum(nil)
	testSig := base64.URLEncoding.EncodeToString(rawSignature)

	if subtle.ConstantTimeCompare([]byte(sig), []byte(testSig)) != 1 {
		c.JSON(401, &errortypes.ErrorData{
			Error:   "authentication_error",
			Message: "Authentication error occurred",
		})
		return
	}

	provider := settings.Auth.GetProvider(tokn.Provider)
	if provider == nil {
		c.AbortWithStatus(404)
		return
	}

	err = tokn.Remove(db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	roles := []string{}
	roles = append(roles, provider.DefaultRoles...)

	for _, role := range strings.Split(params.Get("roles"), ",") {
		if role != "" {
			roles = append(roles, role)
		}
	}

	username := params.Get("username")

	var usr *user.User
	if provider.AutoCreate {
		usr = &user.User{
			Type:     provider.Type,
			Username: username,
			Roles:    roles,
		}

		errData, err := usr.Validate(db)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		if errData != nil {
			logrus.WithFields(logrus.Fields{
				"error":     errData.Error,
				"error_msg": errData.Message,
			}).Error("handlers: Single sign on user validate failed")
			return
		}

		err = usr.Upsert(db)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
	} else {
		usr, err = user.GetUsername(db, provider.Type, username)
		if err != nil {
			switch err.(type) {
			case *database.NotFoundError:
				c.JSON(401, &errortypes.ErrorData{
					Error:   "unauthorized",
					Message: "Not authorized",
				})
			default:
				c.AbortWithError(500, err)
			}
			return
		}
	}

	if usr.Disabled || usr.Administrator != "super" {
		c.JSON(401, &errortypes.ErrorData{
			Error:   "unauthorized",
			Message: "Not authorized",
		})
		return
	}

	cook := cookie.New(c)

	_, err = cook.NewSession(db, usr.Id, true)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.Redirect(302, "/")
}
