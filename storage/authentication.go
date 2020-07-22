package storage

import (
	"errors"
	"github.com/chabad360/covey/models"
	"gorm.io/gorm"
	"strconv"
)

// AddUser adds a User to the database.
func AddUser(u models.User) error {
	result := DB.Exec(`INSERT INTO users(username, password_hash) 
	VALUES($1, crypt($2, gen_salt('bf')));`, u.Username, u.Password)
	return result.Error
}

// UpdateUser updates a User in the database.
func UpdateUser(u models.User, uOld models.User) error {
	result := DB.Model(u).Where("(password_hash = crypt(?, password_hash)) = 't'", uOld.Password).Update(
		"password_hash", gorm.Expr("crypt(?, gen_salt('bf'))", u.Password))
	return result.Error
}

// GetUser returns a User from the database.
func GetUser(u models.User) (string, error) {
	var id int
	result := DB.Raw("SELECT id FROM users WHERE username = ? AND (password_hash = crypt(?, password_hash)) = 't'",
		u.Username, u.Password).First(&id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", result.Error
	}

	return strconv.Itoa(id), nil
}
