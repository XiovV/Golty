package db

type User struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password []byte `db:"password"`
}

func (db *DB) CreateUser(user User) error {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	_, err := db.db.ExecContext(ctx, "INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetNumberOfUsers() (int, error) {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	var count int
	err := db.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM users")
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (db *DB) GetUserByUsername(username string) (User, error) {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	var user User

	err := db.db.GetContext(ctx, &user, "SELECT id, username, password FROM users WHERE username = $1", username)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
