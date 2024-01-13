package db

import (
	"github.com/Masterminds/squirrel"
	"github.com/golang-module/carbon/v2"
	"github.com/google/uuid"
	"github.com/newtoallofthis123/noob_social/views"
)

func (pq *PqInstance) CreateUser(req views.CreateUserReq) (string, error) {

	userId := uuid.New()

	query := pq.Builder.Insert("users").Columns("id", "username", "email", "created_at").Values(userId, req.Username, req.Email, carbon.Now()).Suffix("RETURNING \"id\"").RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	toReturn := ""

	err := query.QueryRow().Scan(&toReturn)
	if err != nil {
		return "", err
	}

	return toReturn, nil
}

func (pq *PqInstance) CreateSession(userId string) (string, error) {

	sessionId := uuid.New()

	query := pq.Builder.Insert("sessions").Columns("id", "user_id", "created_at").Values(sessionId, userId, carbon.Now()).Suffix("RETURNING \"id\"").RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	toReturn := ""

	err := query.QueryRow().Scan(&toReturn)
	if err != nil {
		return "", err
	}

	return toReturn, nil
}

func (pq *PqInstance) GetUserByUsername(username string) (views.User, error) {

	query := pq.Builder.Select("*").From("users").Where(squirrel.Eq{"username": username}).RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	user := views.User{}

	userId := ""

	err := query.QueryRow().Scan(&userId, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		return views.User{}, err
	}

	user.Id, err = uuid.Parse(userId)
	if err != nil {
		return views.User{}, err
	}

	return user, nil
}

func (pq *PqInstance) GetUserByEmail(email string) (views.User, error) {

	query := pq.Builder.Select("*").From("users").Where(squirrel.Eq{"email": email}).RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	user := views.User{}

	userId := ""

	err := query.QueryRow().Scan(&userId, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		return views.User{}, err
	}

	user.Id, err = uuid.Parse(userId)
	if err != nil {
		return views.User{}, err
	}

	return user, nil
}

func (pq *PqInstance) CreateOtp(userId string, otp string) (string, error) {

	otpId := uuid.New()

	query := pq.Builder.Insert("otp").Columns("id", "user_id", "otp", "created_at").Values(otpId, userId, otp, carbon.Now()).Suffix("RETURNING \"id\"").RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	toReturn := ""

	err := query.QueryRow().Scan(&toReturn)
	if err != nil {
		return "", err
	}

	return toReturn, nil
}

func (pq *PqInstance) GetOtp(otp_id string) (string, string, error) {

	query := pq.Builder.Select("otp", "user_id").From("otp").Where(squirrel.Eq{"id": otp_id}).RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	otp := ""
	userId := ""

	err := query.QueryRow().Scan(&otp, &userId)
	if err != nil {
		return "", "", err
	}

	return otp, userId, nil
}

func (pq *PqInstance) DeleteOtp(otpId string) error {

	query := pq.Builder.Delete("otp").Where(squirrel.Eq{"id": otpId}).RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	_, err := query.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (pq *PqInstance) GetUserById(userId string) (views.User, error) {

	query := pq.Builder.Select("*").From("users").Where(squirrel.Eq{"id": userId}).RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	user := views.User{}

	toReturn := ""

	err := query.QueryRow().Scan(&toReturn, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		return views.User{}, err
	}

	user.Id, err = uuid.Parse(toReturn)
	if err != nil {
		return views.User{}, err
	}

	return user, nil
}

func (pq *PqInstance) GetSessionById(sessionId string) (views.Session, error) {

	query := pq.Builder.Select("*").From("sessions").Where(squirrel.Eq{"id": sessionId}).RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	session := views.Session{}

	toReturn := ""

	err := query.QueryRow().Scan(&toReturn, &session.UserId, &session.CreatedAt)
	if err != nil {
		return views.Session{}, err
	}

	session.Id, err = uuid.Parse(toReturn)
	if err != nil {
		return views.Session{}, err
	}

	return session, nil
}

func (pq *PqInstance) DeleteSession(sessionId string) error {
	query := pq.Builder.Delete("sessions").Where(squirrel.Eq{"id": sessionId}).RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	_, err := query.Exec()

	if err != nil {
		return err
	}

	return nil
}

func (pq *PqInstance) CreateProfile(req views.CreateProfileReq) (string, error) {
	query := pq.Builder.Insert("profile").Columns("id", "full_name", "profile_pic", "user_id", "bio", "created_at").Values(uuid.New(), req.FullName, req.ProfilePic, req.UserId, req.Bio, carbon.Now()).Suffix("RETURNING \"id\"").RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	toReturn := ""

	if query.QueryRow().Scan(&toReturn) != nil {
		return "", nil
	}

	return toReturn, nil
}

func (pq *PqInstance) DeleteProfile(profileId string) error {
	query := pq.Builder.Delete("profile").Where(squirrel.Eq{"id": profileId}).RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	_, err := query.Exec()

	if err != nil {
		return err
	}

	return nil
}

func (pq *PqInstance) GetProfileByUser(userId string) (views.Profile, error) {
	query := pq.Builder.Select("*").From("profile").Where(squirrel.Eq{"user_id": userId}).RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	profile := views.Profile{}

	profileId := ""

	err := query.QueryRow().Scan(&profileId, &profile.FullName, &profile.ProfilePic, &profile.UserId, &profile.Bio, &profile.CreatedAt)
	if err != nil {
		return views.Profile{}, err
	}

	profile.Id, err = uuid.Parse(profileId)
	if err != nil {
		return views.Profile{}, err
	}

	return profile, nil
}
