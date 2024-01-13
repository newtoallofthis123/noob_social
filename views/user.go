package views

import "github.com/google/uuid"

type CreateUserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type User struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt string    `json:"created_at"`
}

type Session struct {
	Id        uuid.UUID `json:"id"`
	UserId    string    `json:"user_id"`
	CreatedAt string    `json:"created_at"`
}

type CreateProfileReq struct {
	FullName   string `json:"full_name"`
	ProfilePic string `json:"profile_pic"`
	UserId     string `json:"user_id"`
	Bio        string `json:"bio"`
}

type Profile struct {
	Id         uuid.UUID `json:"id"`
	FullName   string    `json:"full_name"`
	ProfilePic string    `json:"profile_pic"`
	UserId     string    `json:"user_id"`
	Bio        string    `json:"bio"`
	CreatedAt  string    `json:"created_at"`
}
