package db

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"github.com/newtoallofthis123/noob_social/utils"
	"github.com/newtoallofthis123/noob_social/views"
)

type Store interface {
	CreateUser(req views.CreateUserReq) (string, error)
	CreateSession(userId string) (string, error)
	GetUserByUsername(username string) (views.User, error)
	GetUserByEmail(email string) (views.User, error)
	CreateOtp(userId string, otp string) (string, error)
	GetOtp(userId string) (string, string, error)
	DeleteOtp(userId string) error
	GetUserById(userId string) (views.User, error)
	GetSessionById(sessionId string) (views.Session, error)
	DeleteSession(sessionId string) error

	CreateContent(req views.CreateContentRequest) (string, error)
	CreatePost(req views.CreatePostStruct) (string, error)
	GetPost(iden string) (views.Post, error)
	GetContent(contentId string) (views.Content, error)
	GetPostsByUser(userId string) ([]views.FullPost, error)
	CreateLike(userId string, postId string) error
	UpdateTotalLikes(postId, exp string) error
	GetLike(userId string, postId string) (views.Like, error)
	DeleteLike(likeId string) error

	CreateProfile(req views.CreateProfileReq) (string, error)
	DeleteProfile(profileId string) error
	GetProfileByUser(userId string) (views.Profile, error)
	GetAllPictures() ([]string, error)
}

type PqInstance struct {
	Db      *sql.DB
	Builder *squirrel.StatementBuilderType
}

// New Returns a new DbInstance.
// This db instance is not tested and may not be connected.
func New() (*PqInstance, error) {

	env := utils.ReadEnv()

	db, err := sql.Open("postgres", env.ConnString)
	if err != nil {
		return nil, err
	}

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)

	return &PqInstance{db, &psql}, err
}

// InitDb Returns a tested and created DbInstance.
// Inherited from New(), this db instance is tested and connected
// and is also pinged to ensure that the connection is still alive.
func InitDb(refresh bool) (*PqInstance, error) {
	db, err := New()
	if err != nil {
		return nil, err
	}

	err = ping(db.Db)
	if err != nil {
		return nil, err
	}

	err = createTables(refresh, db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Pings the database to ensure that the connection is still alive.
func ping(db *sql.DB) error {
	return db.Ping()
}
