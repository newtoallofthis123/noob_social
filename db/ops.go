package db

// Creates the tables if they do not exist.
// Accepts a refresh bool which will drop the tables if true.
// This is not a DBInstance method because it is only used in the
// once through the program
func createTables(refresh bool, db *PqInstance) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		created_at TIMESTAMP NOT NULL
	);

	CREATE TABLE IF NOT EXISTS profile(
		id UUID PRIMARY KEY,
		profile_pic VARCHAR(255),
		user_id UUID REFERENCES users(id),
		bio TEXT,
		created_at TIMESTAMP NOT NULL
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id UUID PRIMARY KEY,
		user_id UUID REFERENCES users(id),
		created_at TIMESTAMP NOT NULL
	);

	CREATE TABLE IF NOT EXISTS contents (
		id UUID PRIMARY KEY,
		body TEXT,
		image VARCHAR(255),
		video VARCHAR(255),
		post_type VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL
	);

	CREATE TABLE IF NOT EXISTS otp (
		id UUID PRIMARY KEY,
		user_id UUID REFERENCES users(id),
		otp VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL
	);

	CREATE TABLE IF NOT EXISTS posts (
		id TEXT PRIMARY KEY,
		author UUID REFERENCES users(id),
		content UUID REFERENCES contents(id),
		comment_to TEXT REFERENCES posts(id),
		created_at TIMESTAMP NOT NULL
	); 

	CREATE TABLE IF NOT EXISTS likes (
		id UUID PRIMARY KEY,
		user_id UUID REFERENCES users(id),
		post_id TEXT REFERENCES posts(id),
		created_at TIMESTAMP NOT NULL
	);
	`

	if refresh {
		query = `
		DROP TABLE IF EXISTS likes;
		DROP TABLE IF EXISTS posts;
		DROP TABLE IF EXISTS contents;
		DROP TABLE IF EXISTS users;
		DROP TABLE IF EXISTS otp;
		DROP TABLE IF EXISTS sessions;
		` + query
	}

	_, err := db.Db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
