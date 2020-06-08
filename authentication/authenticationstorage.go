package authentication

import (
	"context"

	"github.com/chabad360/covey/storage"
)

// AddUser adds a User to the database.
func AddUser(u credentials) error {
	db := storage.GetPool()
	_, err := db.Exec(context.Background(), "INSERT INTO users(username, password_hash) VALUES($1, crypt($2, gen_salt('bf'));", u.Username, u.Password)
	return err
}

// UpdateUser updates a User in the database.
func UpdateUser(u credentials) error {
	db := storage.GetPool()
	_, err := db.Exec(context.Background(), "UPDATE users SET password_hash = crypt($2, gen_salt('bf')) WHERE username = $1 AND (password_hash = crypt($2, password_hash)) = 't';",
		u.Username, u.Password)
	return err
}

// GetUser returns a User ID from the database.
func GetUser(u credentials) (uint32, error) {
	db := storage.GetPool()
	var id uint32
	if err := db.QueryRow(context.Background(), `SELECT id FROM users 
		WHERE username = $1 AND (password_hash = crypt($2, password_hash)) = 't';`,
		u.Username, u.Password).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
