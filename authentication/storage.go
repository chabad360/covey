package authentication

import (
	"errors"
	"github.com/chabad360/covey/models"
	"gorm.io/gorm"
	"strconv"

	"github.com/chabad360/covey/storage"
)

var db *gorm.DB

// AddUser adds a User to the database.
func AddUser(u models.User) error {
	refreshDB()
	result := db.Exec(`INSERT INTO users(username, password_hash) 
	VALUES($1, crypt($2, gen_salt('bf')));`,
		u.Username, u.Password)
	return result.Error
}

// UpdateUser updates a User in the database.
func UpdateUser(u models.User, uOld models.User) error {
	refreshDB()
	result := db.Model(u).Where("(password_hash = crypt(?, password_hash)) = 't'", uOld.Password).Update(
		"password_hash", gorm.Expr("crypt(?, gen_salt('bf'))", u.Password))
	return result.Error
}

// GetUser returns a User from the database.
func GetUser(u models.User) (string, error) {
	refreshDB()
	var id int
	result := db.Table("users").Where("username = ?", u.Username).Where(
		"(password_hash = crypt(?, password_hash)) = 't'", u.Password).Select("id").First(&id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", result.Error
	}

	return strconv.Itoa(id), nil
}

func refreshDB() {
	if db == nil {
		db = storage.DB
	}
}
