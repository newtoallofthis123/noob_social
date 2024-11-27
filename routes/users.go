package routes

import (
	"bytes"
	"fmt"
	"image/png"

	"github.com/anthonynsimon/bild/transform"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/newtoallofthis123/noob_social/email"
	"github.com/newtoallofthis123/noob_social/templates"
	"github.com/newtoallofthis123/noob_social/utils"
	"github.com/newtoallofthis123/noob_social/views"
)

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

	c.SetCookie("session_id", sessionId, 60*60*24*7, "/", "localhost", false, true)

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

func (api *ApiServer) handleUserCustomize(c *gin.Context) {
	bio := c.PostForm("bio")
	fullName := c.PostForm("full_name")
	var finalName = ""
	profilePic, err := c.FormFile("profile_picture")
	if err == nil {
		err := c.SaveUploadedFile(profilePic, utils.FILEPATH+profilePic.Filename)
		if err != nil {
			return
		}

		finalName, err = utils.CheckPicture(profilePic.Filename, true)
		if err != nil {
			c.String(500, err.Error())
			return
		}
	} else {
		existingPic, ok := c.GetPostForm("existing_pic")
		if !ok {
			finalName = ""
		}

		finalName = existingPic
	}

	bannerName := ""
	bannerPic, err := c.FormFile("banner")
	if err == nil {
		err := c.SaveUploadedFile(bannerPic, utils.FILEPATH+bannerPic.Filename)
		if err != nil {
			return
		}

		bannerName, err = utils.CheckPicture(bannerPic.Filename, false)
		if err != nil {
			c.String(500, err.Error())
			return
		}
	} else {
		existingBanner, ok := c.GetPostForm("existing_banner")
		if !ok {
			bannerName = ""
		}

		bannerName = existingBanner
	}

	userId, ok := c.GetPostForm("user_id")
	if !ok {
		c.String(500, "No user id")
		return
	}

	req := views.CreateProfileReq{
		Bio:        bio,
		FullName:   fullName,
		ProfilePic: finalName,
		Banner:     bannerName,
		UserId:     userId,
	}

	// if there is a profile, delete it
	profile, err := api.store.GetProfileByUser(userId)
	if err == nil {
		err = api.store.DeleteProfile(profile.Id.String())
		if err != nil {
			c.String(500, err.Error())
			return
		}
	}

	_, err = api.store.CreateProfile(req)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.Redirect(302, "/")
}

func (api *ApiServer) handleGetUserAvatar(c *gin.Context) {
	userName := c.Param("u")

	user, err := api.store.GetUserByUsername(userName)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	profile, err := api.store.GetProfileByUser(user.Id.String())
	if err != nil {
		c.Redirect(302, "/defaultAvatar/"+user.Username)
	}

	if profile.ProfilePic == "" {
		c.Redirect(302, "/defaultAvatar/"+user.Username)
	}

	imageFile, err := utils.GetImage(profile.ProfilePic)
	if err != nil {
		c.Redirect(302, "/defaultAvatar/"+user.Username)
	}

	c.Data(200, "image/png", imageFile.Bytes())
}

func (api *ApiServer) handleProfilePage(c *gin.Context) {
	loggedUserId, ok := c.Get("user_id")
	if !ok {
		c.Redirect(302, "/login")
		return
	}

	loggedUser, err := api.store.GetUserById(loggedUserId.(string))
	if err != nil {
		c.String(500, err.Error())
		return
	}

	username := c.Param("username")
	user, err := api.store.GetUserByUsername(username)
	if err != nil {
		user, err = api.store.GetUserById(username)
		if err != nil {
			c.String(500, err.Error())
			return
		}
	}

	posts, err := api.store.GetPostsByUser(user.Id.String())
	if err != nil {
		c.String(500, err.Error())
		return
	}

	for i, j := 0, len(posts)-1; i < j; i, j = i+1, j-1 {
		posts[i], posts[j] = posts[j], posts[i]
	}

	likes, err := api.store.GetUserLikes(user.Id.String())
	if err != nil {
		c.String(500, err.Error())
		return
	}

	for i, j := 0, len(likes)-1; i < j; i, j = i+1, j-1 {
		likes[i], likes[j] = likes[j], likes[i]
	}

	profile, err := api.store.GetProfileByUser(user.Id.String())
	if err != nil {
		c.Redirect(302, "/customization")
		return
	}

	follows, err := api.store.DoesUserFollow(loggedUserId.(string), user.Id.String())
	if err != nil {
		c.String(500, err.Error())
		return
	}

	admin := false
	if loggedUser.Id == user.Id {
		admin = true
	}

	templates.AppLayout(fmt.Sprintf("%s's Profile", profile.FullName), loggedUser.Username, templates.ProfilePage(likes, posts, user.Username, profile, follows, admin)).Render(c.Request.Context(), c.Writer)
}

func (api *ApiServer) handleUserBanner(c *gin.Context) {
	username := c.Params.ByName("username")

	user, err := api.store.GetUserByUsername(username)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	profile, err := api.store.GetProfileByUser(user.Id.String())
	if err != nil {
		c.String(500, err.Error())
		return
	}

	imageFile, err := utils.GetImageFile(profile.Banner)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	croppedImage := transform.Resize(imageFile, 1080, 320, transform.Linear)

	var buff bytes.Buffer
	if png.Encode(&buff, croppedImage) != nil {
		c.String(500, "Error encoding image")
	}

	c.Data(200, "image/png", buff.Bytes())
}

func (api *ApiServer) handleFollowUser(c *gin.Context) {
	userId, ok := c.Get("user_id")
	if !ok {
		c.String(500, "No user id")
		return
	}

	followedUserId := c.PostForm("followe")
	fmt.Println(followedUserId)

	err := api.store.CreateFollow(userId.(string), followedUserId)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.String(200, `
<button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-full"
	hx-post="/unfollowUser" hx-target="#follow" hx-swap="outerHTML">
Unfollow
</button>`)
}

func (api *ApiServer) handleUnfollowUser(c *gin.Context) {
	userId, ok := c.Get("user_id")
	if !ok {
		c.String(500, "No user id")
		return
	}

	followedUserId := c.PostForm("followe")
	fmt.Println(followedUserId)

	err := api.store.DeleteFollow(userId.(string), followedUserId)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.String(200, `
<button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-full"
	hx-post="/followUser" hx-target="#follow" hx-swap="outerHTML">
Follow
</button>`)
}

func (api *ApiServer) handleFeedRecommendation(c *gin.Context) {
	userId, ok := c.GetQuery("user_id")
	if !ok {
		c.String(500, "No user id")
		return
	}

	user, err := api.store.GetUserById(userId)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	following, err := api.store.GetUserFollowing(userId)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"user":      user,
		"following": following,
	})
}
