package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/newtoallofthis123/noob_social/templates"
	"github.com/newtoallofthis123/noob_social/utils"
	"github.com/newtoallofthis123/noob_social/views"
	"github.com/shurcooL/github_flavored_markdown"
)

func (api *ApiServer) handleCreatePost(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.Redirect(302, "/login")
		return
	}

	var createContentReq views.CreateContentRequest

	body := c.PostForm("content")
	postType := c.PostForm("post_type")

	createContentReq.Body = body
	createContentReq.PostType = postType

	image, err := c.FormFile("image")
	var finalName string = ""
	if err == nil {
		c.SaveUploadedFile(image, utils.FILEPATH+image.Filename)

		finalName, err = utils.CheckPicture(image.Filename, false)
		if err != nil {
			// again, we want to display the error to the user
			c.String(200, err.Error())
			return
		}
		createContentReq.Image = finalName
	} else {
		fmt.Println(err)
		fmt.Println("Not processing image")
		createContentReq.Image = ""
	}

	// TODO: Video support
	createContentReq.Video = ""

	content, err := api.store.CreateContent(createContentReq)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	var createPostReq views.CreatePostStruct

	createPostReq.UserId = userID.(string)
	createPostReq.Content.Id = content

	var includeBack = false

	// We are just creating a post, not a comment
	// If we were creating a comment, we would set this to the post ID
	commentId := c.PostForm("comment_id")
	if commentId != "" {
		createPostReq.CommentTo = commentId
		includeBack = true
	}

	postIden, err := api.store.CreatePost(createPostReq)
	if err != nil {
		c.String(500, err.Error())
		fmt.Println(err)
		return
	}

	user, err := api.store.GetUserById(userID.(string))
	if err != nil {
		c.String(500, err.Error())
		return
	}
	if includeBack {
		c.Redirect(302, fmt.Sprintf("/%s/post/%s?back=%s", user.Username, postIden, commentId))
	} else {
		c.Redirect(302, fmt.Sprintf("/%s/post/%s", user.Username, postIden))
	}
}

func (api *ApiServer) handlePostPage(c *gin.Context) {
	postIden := c.Params.ByName("iden")
	username := c.Params.ByName("username")

	post, err := api.store.GetPost(postIden)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	content, err := api.store.GetContent(post.Content)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	profile, err := api.store.GetProfileByUser(post.Author)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	// reduce body to 50 chars
	reducedBody := ""
	if len(content.Body) > 50 {
		reducedBody = content.Body[:50] + "..."
	} else {
		reducedBody = content.Body
	}

	name := ""
	if profile.FullName != "" {
		name = profile.FullName
	} else {
		name = username
	}

	isLiked := false
	user_id, ok := c.Get("user_id")
	if !ok {
		isLiked = false
	} else {
		_, err := api.store.GetLike(user_id.(string), postIden)
		if err != nil {
			isLiked = false
		} else {
			isLiked = true
		}
	}

	back := c.Query("back")
	if back == "" {
		back = "/"
	}

	user, err := api.store.GetUserById(user_id.(string))
	if err != nil {
		c.String(500, err.Error())
		return
	}

	comments, err := api.store.GetComments(postIden)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	templates.AppLayout(fmt.Sprintf("%s: \"%s\" on NoobSocial", name, reducedBody), user.Username, templates.PostPage(false, isLiked, username, post, content, profile, comments, back)).Render(c.Request.Context(), c.Writer)
}

func (api *ApiServer) handleGetMdContent(c *gin.Context) {
	body := c.PostForm("body")

	mdText := github_flavored_markdown.Markdown([]byte(body))

	c.String(200, string(mdText))
}

func (api *ApiServer) handleJsonUserPosts(c *gin.Context) {
	username := c.Params.ByName("username")

	user, err := api.store.GetUserByUsername(username)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	userPosts, err := api.store.GetPostsByUser(user.Id.String())
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.JSON(200, userPosts)
}

func (api *ApiServer) handleUserLike(c *gin.Context) {
	userId, ok := c.Get("user_id")
	if !ok {
		c.Redirect(302, "/login")
		return
	}

	postId := c.PostForm("post_id")

	err := api.store.CreateLike(userId.(string), postId)
	if err != nil {
		fmt.Println(userId, postId, err)
		c.String(500, err.Error())
		return
	}

	err = api.store.UpdateTotalLikes(postId, "total_likes + 1")
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.Header("HX-Refresh", "true")
	c.String(200, "ok")
}

func (api *ApiServer) handleUserUnlike(c *gin.Context) {
	userId, ok := c.Get("user_id")
	if !ok {
		c.Redirect(302, "/login")
		return
	}

	postId := c.PostForm("post_id")

	like, err := api.store.GetLike(userId.(string), postId)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	err = api.store.DeleteLike(like.Id)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	err = api.store.UpdateTotalLikes(postId, "total_likes - 1")
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.Header("HX-Refresh", "true")
	c.String(200, "ok")
}
