package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/newtoallofthis123/noob_social/templates"
)

func (api *ApiServer) handleLoginPage(c *gin.Context) {
	templates.AntiProtected("Login", templates.LoginPage()).Render(c.Request.Context(), c.Writer)
}

func (api *ApiServer) handleOtpPage(c *gin.Context) {
	otpId := c.Query("otp_id")
	if otpId == "" {
		c.Redirect(302, "/login")
		return
	}

	templates.AntiProtected("OTP", templates.OtpLogin(otpId)).Render(c.Request.Context(), c.Writer)
}

func (api *ApiServer) handleHomePage(c *gin.Context) {
	templates.Protected("NoobSocial | Home", templates.Home()).Render(c.Request.Context(), c.Writer)
}
