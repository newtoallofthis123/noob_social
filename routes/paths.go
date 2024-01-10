package routes

import (
	"fmt"

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
	sessionId, err := c.Cookie("session_id")
	if err != nil {
		fmt.Println(err)
		c.Redirect(302, "/v")
		return
	}

	session, err := api.store.GetSessionById(sessionId)
	if err != nil {
		fmt.Println(err)
		c.Redirect(302, "/v")
		return
	}

	user, err := api.store.GetUserById(session.UserId)
	if err != nil {
		fmt.Println(err)
		c.Redirect(302, "/v")
		return
	}

	templates.Protected(templates.AppLayout("NoobSocial", user.Username, templates.Home())).Render(c.Request.Context(), c.Writer)
}

func (api *ApiServer) handleSignUpPage(c *gin.Context) {
	templates.AntiProtected("Sign Up", templates.SignUpPage(c.Query("email"))).Render(c.Request.Context(), c.Writer)
}
