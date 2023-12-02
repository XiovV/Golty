package repository

type User struct {
	ID       int
	Username string
	Password string
	Admin    int
}

func (r *Repository) FindUserByUsername(username string) (User, error) {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	var user User
	err := r.db.GetContext(ctx, &user, "SELECT id, username, password FROM users WHERE username = $1", username)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *Repository) InsertUser(user User) error {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	_, err := r.db.ExecContext(ctx, "INSERT INTO users (username, password, admin) VALUES ($1, $2, $3)", user.Username, user.Password, user.Admin)
	if err != nil {
		return err
	}

	return nil
}
