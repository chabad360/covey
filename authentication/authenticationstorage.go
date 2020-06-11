package authentication

import (
	"context"
	"strconv"

	"github.com/chabad360/covey/storage"
	"github.com/jackc/pgx/v4/pgxpool"
)

var db *pgxpool.Pool

// AddUser adds a User to the database.
func AddUser(u credentials) error {
	refreshDB()
	_, err := db.Exec(context.Background(), `INSERT INTO users(username, password_hash) 
	VALUES($1, crypt($2, gen_salt('bf')));`,
		u.Username, u.Password)
	return err
}

// UpdateUser updates a User in the database.
func UpdateUser(u credentials, uOld credentials) error {
	refreshDB()
	_, err := db.Exec(context.Background(), `UPDATE users SET password_hash = crypt($2, gen_salt('bf')) 
		WHERE username = $1 AND (password_hash = crypt($3, password_hash)) = 't';`,
		u.Username, u.Password, uOld.Password)
	return err
}

// GetUser returns a User ID from the database.
func GetUser(u credentials) (string, error) {
	refreshDB()
	var id int
	if err := db.QueryRow(context.Background(), `SELECT id FROM users 
		WHERE username = $1 AND (password_hash = crypt($2, password_hash)) = 't';`,
		u.Username, u.Password).Scan(&id); err != nil {
		return "", err
	}
	return strconv.Itoa(id), nil
}

func refreshDB() {
	if db == nil {
		db = storage.DB
	}
}
