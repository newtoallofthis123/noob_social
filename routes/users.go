package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/newtoallofthis123/noob_social/email"
	"github.com/newtoallofthis123/noob_social/utils"
	"github.com/newtoallofthis123/noob_social/views"
)

func (api *ApiServer) handleCreateUser(c *gin.Context) {
	test := views.CreateUserReq{
		Username: "noob",
		Email:    "noob@duck.com",
	}

	user, err := api.store.CreateUser(test)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"user": user,
	})
}

func (api *ApiServer) handleSendOTP(c *gin.Context) {
	to := c.Query("to")

	err := email.SendOtp(utils.GenerateOtp(8), "test", to)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Sent OTP",
	})
}

func (api *ApiServer) handleEmailLogin(c *gin.Context) {
	userEmail := c.PostForm("email")

	user, err := api.store.GetUserByEmail(userEmail)
	if err != nil {
		c.Redirect(302, "/sign-up?email="+userEmail)
		return
	}

	otp := utils.GenerateOtp(8)

	err = email.SendOtp(otp, user.Username, user.Email)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	optId, err := api.store.CreateOtp(user.Id.String(), otp)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.Redirect(302, "/otp-login?otp_id="+optId)
}

func (api *ApiServer) handleOtpLogin(c *gin.Context) {
	otpId := c.PostForm("otp_id")
	submittedOtp := c.PostForm("otp")

	otp, userId, err := api.store.GetOtp(otpId)
	if err != nil {
		fmt.Println(err)
		c.String(500, err.Error())
		return
	}

	if otp != submittedOtp {
		fmt.Println(err)
		c.String(500, "Invalid OTP")
		return
	}

	sessionId, err := api.store.CreateSession(userId)
	if err != nil {
		fmt.Println(err)
		c.String(500, err.Error())
		return
	}

	c.SetCookie("session_id", sessionId, 3600, "/", "localhost", false, true)

	err = api.store.DeleteOtp(otpId)
	if err != nil {
		fmt.Println(err)
		c.String(500, err.Error())
		return
	}

	c.Redirect(302, "/")
}

func (api *ApiServer) handleSignOut(c *gin.Context) {

	cookie, err := c.Cookie("session_id")
	if err != nil {
		c.String(500, err.Error())
		return
	}

	// Delete a cookie by setting the max age to -1
	// Hence effectively deleting it
	c.SetCookie("session_id", "", -1, "/", "localhost", false, true)

	// TODO: Let something actually happen if this fails
	_ = api.store.DeleteSession(cookie)

	c.String(200, "Signed out")
}

func (api *ApiServer) handleUserSignUp(c *gin.Context) {
	username := c.PostForm("username")
	userEmail := c.PostForm("email")

	userId, err := api.store.CreateUser(views.CreateUserReq{
		Username: username,
		Email:    userEmail,
	})
	if err != nil {
		c.String(500, err.Error())
		return
	}

	parsed, err := uuid.Parse(userId)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	user := views.User{
		Id:       parsed,
		Username: username,
		Email:    userEmail,
	}

	// Instead of queries, we can just create a mock user
	otp := utils.GenerateOtp(8)

	err = email.SendOtp(otp, user.Username, user.Email)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	optId, err := api.store.CreateOtp(user.Id.String(), otp)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.Redirect(302, "/otp-login?otp_id="+optId)
}
