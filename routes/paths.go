package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/newtoallofthis123/noob_social/templates"
	"github.com/newtoallofthis123/noob_social/views"
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
		c.Redirect(302, "/login")
		return
	}

	// // FIX: THIS IS SO BAD! BUT I MEAN I DON'T KNOW HOW TO DO IT ANY OTHER WAY!
	// // This is used to delete unused images
	// randomNum := rand.Intn(10)
	// if randomNum == 6 {
	// 	usedImage, err := api.store.GetAllPictures()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		c.Redirect(302, "/err")
	// 		return
	// 	}
	//
	// 	if utils.DeleteUnused(usedImage) != nil {
	// 		fmt.Println(err)
	// 	}
	// }

	session, err := api.store.GetSessionById(sessionId)
	if err != nil {
		fmt.Println(err)
		c.Redirect(302, "/logout")
		return
	}

	user, err := api.store.GetUserById(session.UserId)
	if err != nil {
		fmt.Println(err)
		c.Redirect(302, "/err")
		return
	}

	templates.Protected(templates.AppLayout("NoobSocial", user.Username, templates.Home())).Render(c.Request.Context(), c.Writer)
}

func (api *ApiServer) handleSignUpPage(c *gin.Context) {
	templates.AntiProtected("Sign Up", templates.SignUpPage(c.Query("email"))).Render(c.Request.Context(), c.Writer)
}

func (api *ApiServer) handleCustomizationPage(c *gin.Context) {
	toUpdate := false
	isNew := c.Query("new")
	if isNew == "" {
		toUpdate = true
	}

	userId, ok := c.Get("user_id")
	fmt.Println(userId)
	if !ok {
		c.Redirect(302, "/login")
		return
	}

	profile, err := api.store.GetProfileByUser(userId.(string))
	if err == nil {
		templates.Protected(templates.Base("NoobSocial", templates.CustomizationPage(userId.(string), toUpdate, profile))).Render(c.Request.Context(), c.Writer)
		return
	}

	userIdStr := userId.(string)

	templates.Protected(templates.Base("NoobSocial", templates.CustomizationPage(userIdStr, toUpdate, views.Profile{}))).Render(c.Request.Context(), c.Writer)
}
