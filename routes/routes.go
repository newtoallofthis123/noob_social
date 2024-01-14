package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/newtoallofthis123/noob_social/db"
	"github.com/newtoallofthis123/noob_social/templates"
	"github.com/newtoallofthis123/noob_social/utils"
)

type ApiServer struct {
	listenAddr string
	store      db.Store
}

// Returns a new Tested ApiServer.
func New() *ApiServer {
	listenAddr := utils.ReadEnv().ListenAddr
	config := utils.ReadConfig()
	store, err := db.InitDb(config.DbRefresh)
	if err != nil {
		panic(err)
	}

	return &ApiServer{listenAddr, store}
}

// Starts the server.
func (api *ApiServer) Start() error {
	r := gin.Default()

	// There are only two static routes, /public and /static.
	r.Static("/public", "./public")
	r.Static("/static", "./static")

	// The favicon
	r.StaticFile("/favicon.ico", "./public/favicon.ico")

	// A small test handler
	r.GET("/version", func(c *gin.Context) {
		templates.Base("Version Route", templates.TestRoute("v.0.0.1")).Render(c.Request.Context(), c.Writer)
	})

	// The home page
	r.GET("/", api.handleHomePage)

	// Test handler for creating a user
	r.GET("/create-user", api.handleCreateUser)
	// Test handler for sending an OTP
	r.GET("/send-otp", api.handleSendOTP)

	r.GET("/login", api.handleLoginPage)
	r.GET("/otp-login", api.handleOtpPage)
	r.GET("/sign-up", api.handleSignUpPage)

	r.POST("/loginUser", api.handleEmailLogin)
	r.POST("/checkOtp", api.handleOtpLogin)
	r.POST("/signUpUser", api.handleUserSignUp)

	r.GET("/checkSession", api.handleAuthCheck)
	r.POST("/logout", api.handleSignOut)

	// util routes
	r.GET("/defaultAvatar/:u", handleDefaultAvatar)
	r.GET("/avatar/:u", api.handleGetUserAvatar)
	r.POST("/md", api.handleGetMdContent)

	// Some JSON routes for testing
	r.GET("/json/:username/posts", api.handleJsonUserPosts)

	// authenticated routes
	auth := r.Group("/")
	auth.Use(api.authMiddleware())

	auth.POST("/createPost", api.handleCreatePost)
	auth.GET("/customization", api.handleCustomizationPage)
	auth.POST("/customizeUser", api.handleUserCustomize)

	auth.POST("/likePost", api.handleUserLike)
	auth.POST("/unlikePost", api.handleUserUnlike)

	auth.GET("/:username/post/:iden", api.handlePostPage)
	auth.GET("/:username", api.handleProfilePage)

	err := r.Run(api.listenAddr)
	return err
}
