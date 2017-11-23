package uhandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/authorizer"
)

func sshValidateGet(c *gin.Context) {
	authr := c.MustGet("authorizer").(*authorizer.Authorizer)

	if !authr.IsValid() {
		if c.Request.URL.RawQuery == "" {
			c.Redirect(302, "/login")
		} else {
			query := c.Request.URL.Query()
			query.Set("redirect", "ssh-validate")
			c.Redirect(302, "/login?"+query.Encode())
		}
		return
	}

	c.String(200, "validated")
}
