package uhandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/audit"
	"github.com/pritunl/pritunl-zero/authorizer"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/sshcert"
	"github.com/pritunl/pritunl-zero/utils"
	"time"
)

type sshValidateData struct {
	Token     string `json:"token"`
	PublicKey string `json:"public_key,omitempty"`
}

type sshCertificateData struct {
	Token        string   `json:"token"`
	Certificates []string `json:"certificates"`
}

func sshGet(c *gin.Context) {
	authr := c.MustGet("authorizer").(*authorizer.Authorizer)

	redirect := ""

	if authr.IsValid() {
		if c.Request.URL.RawQuery == "" {
			redirect = "/"
		} else {
			query := c.Request.URL.Query()
			redirect = "/?" + query.Encode()
		}
	} else {
		if c.Request.URL.RawQuery == "" {
			redirect = "/login"
		} else {
			query := c.Request.URL.Query()
			redirect = "/login?" + query.Encode()
		}
	}

	c.Redirect(302, redirect)
}

func sshValidatePut(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	authr := c.MustGet("authorizer").(*authorizer.Authorizer)

	sshToken := c.Param("ssh_token")

	usr, err := authr.GetUser(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	chal, err := sshcert.GetChallenge(db, sshToken)
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

	err = audit.New(
		db,
		c.Request,
		usr.Id,
		audit.SshAprove,
		audit.Fields{
			"ssh_key": chal.PubKey,
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	err = chal.Approve(db, usr, c.Request)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.Publish(db, "ssh_challenge", chal.Id)

	c.Status(200)
}

func sshValidateDelete(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	authr := c.MustGet("authorizer").(*authorizer.Authorizer)

	sshToken := c.Param("ssh_token")

	usr, err := authr.GetUser(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	chal, err := sshcert.GetChallenge(db, sshToken)
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

	err = audit.New(
		db,
		c.Request,
		usr.Id,
		audit.SshDeny,
		audit.Fields{
			"ssh_key": chal.PubKey,
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	err = chal.Deny(db, usr)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.Publish(db, "ssh_challenge", chal.Id)

	c.Status(200)
}

func sshChallengePut(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := &sshValidateData{}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	chal, err := sshcert.GetChallenge(db, data.Token)
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
	token := chal.Id

	sync := func() {
		chal, err = sshcert.GetChallenge(db, data.Token)
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
		switch chal.State {
		case sshcert.Approved:
			cert, err := sshcert.GetCertificate(db, chal.CertificateId)
			if err != nil {
				switch err.(type) {
				case *database.NotFoundError:
					utils.AbortWithStatus(c, 404)
					break
				default:
					utils.AbortWithError(c, 500, err)
				}
				return true
			}

			resp := &sshCertificateData{
				Token:        token,
				Certificates: cert.Certificates,
			}

			c.JSON(200, resp)

			return true
		case sshcert.Unavailable:
			c.Status(412)
			return true
		case sshcert.Denied:
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

	listenerId := sshcert.Register(token, func() {
		defer func() {
			recover()
		}()
		notify <- true
	})
	defer sshcert.Unregister(token, listenerId)

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

func sshChallengePost(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := &sshValidateData{}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	chal, err := sshcert.NewChallenge(db, data.PublicKey)
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

	resp := &sshValidateData{
		Token: chal.Id,
	}

	c.JSON(200, resp)
}
