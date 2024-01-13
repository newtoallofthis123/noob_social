package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	htmlmaker "github.com/newtoallofthis123/html_maker"
	"github.com/newtoallofthis123/noob_social/templates"
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

	if postType == "img" {
		// TODO: Upload image to local disk
		createContentReq.Image = ""
	} else if postType == "vid" {
		// TODO: Upload video to local disk
		createContentReq.Video = ""
	} else {
		createContentReq.Image = ""
		createContentReq.Video = ""
	}

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

	// We are just creating a post, not a comment
	// If we were creating a comment, we would set this to the post ID
	createPostReq.CommentTo = ""

	postIden, err := api.store.CreatePost(createPostReq)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := api.store.GetUserById(userID.(string))
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	tag := htmlmaker.New("div")

	linkTag := htmlmaker.New("a")
	linkTag.AddAttr("href", fmt.Sprintf("/%s/post/%s", user.Username, postIden))
	linkTag.AddClasses([]string{"underline", "text-blue-500", "hover:text-blue-700"})
	linkTag.Body = "Created post with ID " + postIden

	tag.AddChild(linkTag)

	c.String(200, tag.Convert())
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

	templates.AppLayout(fmt.Sprintf("%s: \"%s\" on NoobSocial", name, reducedBody), username, templates.PostPage(username, post, content)).Render(c.Request.Context(), c.Writer)
}

func (api *ApiServer) handleGetMdContent(c *gin.Context) {
	body := c.PostForm("body")

	mdText := github_flavored_markdown.Markdown([]byte(body))

	fmt.Println(string(mdText))

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
