package authentication

import (
	"strconv"

	"github.com/chabad360/covey/storage"
)

// AddUser adds a User to the database.
func AddUser(u credentials) error {
	db := storage.GetDB()
	_, err := db.Exec("INSERT INTO users(username, password_hash) VALUES($1, crypt($2, gen_salt('bf'));", u.Username, u.Password)
	return err
}

// UpdateUser updates a User in the database.
func UpdateUser(u credentials) error {
	db := storage.GetDB()
	_, err := db.Exec("UPDATE users SET password_hash = crypt($2, gen_salt('bf')) WHERE username = $1 AND (password_hash = crypt($2, password_hash)) = 't';",
		u.Username, u.Password)
	return err
}

// GetUser returns a User ID from the database.
func GetUser(u credentials) (string, error) {
	db := storage.GetDB()
	var id int
	if err := db.QueryRow(`SELECT id FROM users 
		WHERE username = $1 AND (password_hash = crypt($2, password_hash)) = 't';`,
		u.Username, u.Password).Scan(&id); err != nil {
		return "", err
	}
	return strconv.Itoa(id), nil
}
