package uhandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/audit"
	"github.com/pritunl/pritunl-zero/authorizer"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/keybase"
	"github.com/pritunl/pritunl-zero/ssh"
	"github.com/pritunl/pritunl-zero/user"
	"github.com/pritunl/pritunl-zero/utils"
	"time"
)

type keybaseAssociateData struct {
	Username string `json:"username"`
}

type keybaseValidateData struct {
	Token     string `json:"token"`
	Message   string `json:"message,omitempty"`
	Signature string `json:"signature,omitempty"`
}

func keybaseValidatePut(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	authr := c.MustGet("authorizer").(*authorizer.Authorizer)
	data := &keybaseValidateData{}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	usr, err := authr.GetUser(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	asc, err := keybase.GetAssociation(db, data.Token)
	if err != nil {
		switch err.(type) {
		case *database.NotFoundError:
			utils.AbortWithStatus(c, 404)
			break
		default:
			utils.AbortWithError(c, 500, err)
		}
		return
	}

	err, errData := asc.Validate(data.Signature)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(406, errData)
		return
	}

	err = audit.New(
		db,
		c.Request,
		usr.Id,
		audit.KeybaseAssociationAprove,
		audit.Fields{
			"keybase_username": asc.Username,
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	err, errData = asc.Approve(db, usr)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.Publish(db, "keybase_association", asc.Id)

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	c.Status(200)
}

func keybaseValidateDelete(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	authr := c.MustGet("authorizer").(*authorizer.Authorizer)
	data := &keybaseValidateData{}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	usr, err := authr.GetUser(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	asc, err := keybase.GetAssociation(db, data.Token)
	if err != nil {
		switch err.(type) {
		case *database.NotFoundError:
			utils.AbortWithStatus(c, 404)
			break
		default:
			utils.AbortWithError(c, 500, err)
		}
		return
	}

	err, errData := asc.Validate(data.Signature)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = audit.New(
		db,
		c.Request,
		usr.Id,
		audit.KeybaseAssociationDeny,
		audit.Fields{
			"keybase_username": asc.Username,
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	err = asc.Deny(db, usr)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.Publish(db, "keybase_association", asc.Id)

	c.Status(200)
}

func keybaseCheckPut(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := &keybaseValidateData{}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	asc, err := keybase.GetAssociation(db, data.Token)
	if err != nil {
		switch err.(type) {
		case *database.NotFoundError:
			utils.AbortWithStatus(c, 404)
			break
		default:
			utils.AbortWithError(c, 500, err)
		}
		return
	}

	err, errData := asc.Validate(data.Signature)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	_, err = user.GetKeybase(db, asc.Username)
	if err != nil {
		switch err.(type) {
		case *database.NotFoundError:
			utils.AbortWithStatus(c, 404)
			break
		default:
			utils.AbortWithError(c, 500, err)
		}
		return
	}

	c.Status(200)
}

func keybaseAssociatePost(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := &keybaseAssociateData{}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	asc, err := keybase.NewAssociation(db, data.Username)
	if err != nil {
		switch err.(type) {
		case *database.NotFoundError:
			utils.AbortWithStatus(c, 404)
			break
		default:
			utils.AbortWithError(c, 500, err)
		}
		return
	}

	resp := &keybaseValidateData{
		Token:   asc.Id,
		Message: asc.Message(),
	}

	c.JSON(200, resp)
}

func keybaseAssociateGet(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	token := c.Param("token")

	asc, err := keybase.GetAssociation(db, token)
	if err != nil {
		switch err.(type) {
		case *database.NotFoundError:
			utils.AbortWithStatus(c, 404)
			break
		default:
			utils.AbortWithError(c, 500, err)
		}
		return
	}
	token = asc.Id

	sync := func() {
		asc, err = keybase.GetAssociation(db, token)
		if err != nil {
			switch err.(type) {
			case *database.NotFoundError:
				utils.AbortWithStatus(c, 404)
				break
			default:
				utils.AbortWithError(c, 500, err)
			}
			return
		}
	}

	update := func() bool {
		switch asc.State {
		case keybase.Approved:
			c.Status(200)
			return true
		case keybase.Denied:
			c.Status(401)
			return true
		}

		return false
	}

	if update() {
		return
	}

	start := time.Now()
	ticker := time.NewTicker(3 * time.Second)
	notify := make(chan bool, 3)

	listenerId := keybase.Register(token, func() {
		defer func() {
			recover()
		}()
		notify <- true
	})
	defer keybase.Unregister(token, listenerId)

	for {
		select {
		case <-ticker.C:
			if time.Since(start) > 29*time.Second {
				c.Status(205)
				return
			}

			sync()
			if update() {
				return
			}
		case <-notify:
			sync()
			if update() {
				return
			}
		}
	}
}

type keybaseChallengeData struct {
	Username  string `json:"username"`
	PublicKey string `json:"public_key"`
}

type keybaseChallengeRespData struct {
	Token     string `json:"token"`
	Message   string `json:"message"`
	Signature string `json:"signature,omitempty"`
}

type keybaseCertificateData struct {
	Token        string   `json:"token"`
	Certificates []string `json:"certificates"`
}

func keybaseChallengePost(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := &keybaseChallengeData{}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	chal, err := keybase.NewChallenge(db, data.Username, data.PublicKey)
	if err != nil {
		switch err.(type) {
		case *database.NotFoundError:
			utils.AbortWithStatus(c, 404)
			break
		default:
			utils.AbortWithError(c, 500, err)
		}
		return
	}

	resp := &keybaseValidateData{
		Token:   chal.Id,
		Message: chal.Message(),
	}

	c.JSON(200, resp)
}

func keybaseChallengePut(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := &keybaseChallengeRespData{}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	chal, err := keybase.GetChallenge(db, data.Token)
	if err != nil {
		switch err.(type) {
		case *database.NotFoundError:
			utils.AbortWithStatus(c, 404)
			break
		default:
			utils.AbortWithError(c, 500, err)
		}
		return
	}

	cert, err, errData := chal.Validate(db, c.Request, data.Signature)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(406, errData)
		return
	}

	resp := &keybaseCertificateData{
		Token:        chal.Id,
		Certificates: cert.Certificates,
	}

	c.JSON(200, resp)
}
