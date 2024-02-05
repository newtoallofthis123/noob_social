package db

import (
	"github.com/Masterminds/squirrel"
	"github.com/golang-module/carbon/v2"
	"github.com/google/uuid"
	"github.com/newtoallofthis123/noob_social/utils"
	"github.com/newtoallofthis123/noob_social/views"
)

// TODO: FIX Time to Local Time
func (pq *PqInstance) CreateContent(req views.CreateContentRequest) (string, error) {
	contentId := uuid.New()

	query := pq.Builder.Insert("contents").Columns("id", "body", "image", "video", "post_type", "created_at").Values(contentId, req.Body, req.Image, req.Video, req.PostType, carbon.Now(carbon.Local)).Suffix("RETURNING \"id\"").RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	toReturn := ""

	err := query.QueryRow().Scan(&toReturn)
	if err != nil {
		return "", err
	}

	return toReturn, nil
}

func (pq *PqInstance) CreatePost(req views.CreatePostStruct) (string, error) {
	postId := utils.GenerateRandomString(18)

	query := pq.Builder.Insert("posts").Columns("id", "author", "content", "comment_to", "created_at").Values(postId, req.UserId, req.Content.Id, req.CommentTo, carbon.Now(carbon.Local)).Suffix("RETURNING \"id\"").RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	toReturn := ""

	err := query.QueryRow().Scan(&toReturn)
	if err != nil {
		return "", err
	}

	return toReturn, nil
}

func (pq *PqInstance) GetPost(iden string) (views.Post, error) {
	query := pq.Builder.Select("*").From("posts").Where(squirrel.Eq{"id": iden}).RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	post := views.Post{}

	err := query.QueryRow().Scan(&post.Id, &post.Author, &post.Content, &post.TotalLikes, &post.CommentTo, &post.CreatedAt)
	if err != nil {
		return views.Post{}, err
	}

	return post, nil
}

func (pq *PqInstance) GetContent(contentId string) (views.Content, error) {
	query := pq.Builder.Select("*").From("contents").Where(squirrel.Eq{"id": contentId}).RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	content := views.Content{}

	err := query.QueryRow().Scan(&content.Id, &content.Body, &content.Image, &content.Video, &content.PostType, &content.CreatedAt)
	if err != nil {
		return views.Content{}, err
	}

	return content, nil
}

func (pq *PqInstance) GetPostsByUser(userId string) ([]views.FullPost, error) {
	query := pq.Builder.Select("*").From("posts").Where(squirrel.Eq{"author": userId}).RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	var posts []views.FullPost

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post views.FullPost
		err := rows.Scan(&post.Post.Id, &post.Post.Author, &post.Post.Content, &post.Post.TotalLikes, &post.Post.CommentTo, &post.Post.CreatedAt)
		if err != nil {
			return nil, err
		}

		content, err := pq.GetContent(post.Post.Content)
		if err != nil {
			return nil, err
		}

		post.Content = content

		posts = append(posts, post)
	}

	return posts, nil
}

func (pq *PqInstance) CreateLike(userId, postId string) error {
	likeId := uuid.New()

	query := pq.Builder.Insert("likes").Columns("id", "user_id", "post_id", "created_at").Values(likeId, userId, postId, carbon.Now(carbon.Local)).RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	_, err := query.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (pq *PqInstance) DeleteLike(likeId string) error {
	query := pq.Builder.Delete("likes").Where(squirrel.Eq{"id": likeId}).RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	_, err := query.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (pq *PqInstance) GetLike(userId, postId string) (views.Like, error) {
	query := pq.Builder.Select("*").From("likes").Where(squirrel.Eq{"user_id": userId, "post_id": postId}).RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	like := views.Like{}

	err := query.QueryRow().Scan(&like.Id, &like.UserId, &like.PostId, &like.CreatedAt)
	if err != nil {
		return views.Like{}, err
	}

	return like, nil
}

func (pq *PqInstance) UpdateTotalLikes(postId, exp string) error {
	query := pq.Builder.Update("posts").Set("total_likes", squirrel.Expr(exp)).Where(squirrel.Eq{"id": postId}).RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	_, err := query.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (pq *PqInstance) GetComments(postId string) ([]views.Comment, error) {
	query := pq.Builder.Select("posts.id", "posts.author", "posts.content", "posts.total_likes", "posts.comment_to", "posts.created_at", "contents.id", "contents.body", "contents.image", "contents.video", "contents.post_type", "contents.created_at", "users.username", "profile.id", "profile.full_name", "profile.bio", "profile.profile_pic").From("posts").Join("contents ON posts.content = contents.id").Join("users ON posts.author = users.id").Join("profile ON users.id = profile.user_id").Where(squirrel.Eq{"posts.comment_to": postId}).RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	var comments []views.Comment

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var comment views.Comment
		err := rows.Scan(&comment.Post.Id, &comment.Post.Author, &comment.Post.Content, &comment.Post.TotalLikes, &comment.Post.CommentTo, &comment.Post.CreatedAt, &comment.Content.Id, &comment.Content.Body, &comment.Content.Image, &comment.Content.Video, &comment.Content.PostType, &comment.Content.CreatedAt, &comment.Username, &comment.Profile.Id, &comment.Profile.FullName, &comment.Profile.Bio, &comment.Profile.ProfilePic)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	return comments, nil
}
