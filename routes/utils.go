package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/newtoallofthis123/noob_social/utils"
)

func handleDefaultAvatar(c *gin.Context) {
	username := c.Param("u")
	if username == "" {
		c.Redirect(302, "/err")
		return
	}

	avatar, err := utils.GetAvatar(username)
	if err != nil {
		c.Redirect(302, "/err")
		return
	}

	c.Data(200, "image/png", avatar.Bytes())
}
