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
	query := pq.Builder.Insert("profile").Columns("id", "full_name", "profile_pic", "banner", "user_id", "bio", "created_at").Values(uuid.New(), req.FullName, req.ProfilePic, req.Banner, req.UserId, req.Bio, carbon.Now()).Suffix("RETURNING \"id\"").RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

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

	err := query.QueryRow().Scan(&profileId, &profile.FullName, &profile.ProfilePic, &profile.Banner, &profile.UserId, &profile.Bio, &profile.CreatedAt)
	if err != nil {
		return views.Profile{}, err
	}

	profile.Id, err = uuid.Parse(profileId)
	if err != nil {
		return views.Profile{}, err
	}

	return profile, nil
}

func (pq *PqInstance) GetAllPictures() ([]string, error) {
	query := pq.Builder.Select("profile_pic", "banner").From("profile").RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	rows, err := query.Query()
	if err != nil {
		return []string{}, err
	}

	toReturn := []string{}

	for rows.Next() {
		var profilePic string
		var banner string

		err := rows.Scan(&profilePic, &banner)
		if err != nil {
			return []string{}, err
		}

		toReturn = append(toReturn, profilePic, banner)
	}

	return toReturn, nil
}

func (pq *PqInstance) GetUserLikes(userId string) ([]views.Like, error) {
	query := pq.Builder.Select("*").From("likes").Where(squirrel.Eq{"user_id": userId}).RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	rows, err := query.Query()
	if err != nil {
		return []views.Like{}, err
	}

	toReturn := []views.Like{}

	for rows.Next() {
		like := views.Like{}

		err := rows.Scan(&like.Id, &like.UserId, &like.PostId, &like.CreatedAt)
		if err != nil {
			return []views.Like{}, err
		}

		toReturn = append(toReturn, like)
	}

	return toReturn, nil
}

func (pq *PqInstance) GetUserFollowing(userId string) ([]views.User, error) {
	query := pq.Builder.Select("users.id", "users.username", "users.email", "users.created_at").From("users").Join("follows ON users.id = follows.following_id").Where(squirrel.Eq{"follows.user_id": userId}).RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	rows, err := query.Query()
	if err != nil {
		return []views.User{}, err
	}

	toReturn := []views.User{}

	for rows.Next() {
		user := views.User{}

		userId := ""

		err := rows.Scan(&userId, &user.Username, &user.Email, &user.CreatedAt)
		if err != nil {
			return []views.User{}, err
		}

		user.Id, err = uuid.Parse(userId)
		if err != nil {
			return []views.User{}, err
		}

		toReturn = append(toReturn, user)
	}

	return toReturn, nil
}

func (pq *PqInstance) CreateFollow(userId, followId string) error {
	query := pq.Builder.Insert("follows").Columns("id", "user_id", "following_id", "created_at").Values(uuid.New(), userId, followId, carbon.Now()).RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	_, err := query.Exec()

	if err != nil {
		return err
	}

	return nil
}

func (pq *PqInstance) DeleteFollow(userId, followId string) error {
	query := pq.Builder.Delete("follows").Where(squirrel.Eq{"user_id": userId, "following_id": followId}).RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	_, err := query.Exec()

	if err != nil {
		return err
	}

	return nil
}

func (pq *PqInstance) DoesUserFollow(userId, followId string) (bool, error) {
	query := pq.Builder.Select("*").From("follows").Where(squirrel.Eq{"user_id": userId, "following_id": followId}).RunWith(pq.Db).PlaceholderFormat(squirrel.Dollar)

	rows, err := query.Query()
	if err != nil {
		return false, err
	}

	return rows.Next(), nil
}
