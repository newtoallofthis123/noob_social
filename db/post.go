package db

import (
	"github.com/Masterminds/squirrel"
	"github.com/golang-module/carbon/v2"
	"github.com/google/uuid"
	"github.com/newtoallofthis123/noob_social/utils"
	"github.com/newtoallofthis123/noob_social/views"
)

func (pq *PqInstance) CreateContent(req views.CreateContentRequest) (string, error) {
	contentId := uuid.New()

	query := pq.Builder.Insert("contents").Columns("id", "body", "image", "video", "post_type", "created_at").Values(contentId, req.Body, req.Image, req.Video, req.PostType, carbon.Now()).Suffix("RETURNING \"id\"").RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	toReturn := ""

	err := query.QueryRow().Scan(&toReturn)
	if err != nil {
		return "", err
	}

	return toReturn, nil
}

func (pq *PqInstance) CreatePost(req views.CreatePostStruct) (string, error) {
	postId := utils.GenerateRandomString(18)

	query := pq.Builder.Insert("posts").Columns("id", "author", "content", "created_at").Values(postId, req.UserId, req.Content.Id, carbon.Now()).Suffix("RETURNING \"id\"").RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

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

	commentId := uuid.Nil

	err := query.QueryRow().Scan(&post.Id, &post.Author, &post.Content, &commentId, &post.CreatedAt)
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

	var commentId uuid.UUID

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post views.FullPost
		err := rows.Scan(&post.Post.Id, &post.Post.Author, &post.Post.Content, &commentId, &post.Post.CreatedAt)
		if err != nil {
			return nil, err
		}

		if commentId != uuid.Nil {
			post.Post.CommentTo = commentId.String()
		} else {
			post.Post.CommentTo = ""
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
