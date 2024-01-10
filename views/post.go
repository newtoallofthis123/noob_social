package views

// CreatePostStruct is a struct that contains the request body for creating a post
type CreatePostStruct struct {
	UserId    string
	Content   Content
	CommentTo string
}

// CreateContentRequest is a struct that contains the request body for the content
// of a post
type CreateContentRequest struct {
	Body     string
	Image    string
	Video    string
	PostType string
}

// Content is a struct that contains the content of a post
type Content struct {
	Id        string `json:"id"`
	Body      string `json:"body"`
	Image     string `json:"image"`
	Video     string `json:"video"`
	PostType  string `json:"post_type"`
	CreatedAt string `json:"created_at"`
}

// Post is a struct that contains the post
type Post struct {
	Id        string `json:"id"`
	Author    string `json:"author"`
	Content   string `json:"content"`
	CommentTo string `json:"comment_to"`
	CreatedAt string `json:"created_at"`
}
