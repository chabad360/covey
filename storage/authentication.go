package storage

import (
	"github.com/chabad360/covey/models"
	"gorm.io/gorm"
	"log"
)

// AddUser adds a User to the database.
func AddUser(u models.User) error {
	result := DB.Exec(`INSERT INTO users(username, password_hash) 
	VALUES(?, crypt(?, gen_salt('bf')));`, u.Username, u.Password)
	return result.Error
}

// UpdateUser updates a User in the database.
func UpdateUser(u models.User, uOld models.User) error {
	result := DB.Table("users").Where("username = ?", u.Username).Where("(password_hash = crypt(?, password_hash)) = 't'", uOld.Password).Update(
		"password_hash", gorm.Expr("crypt(?, gen_salt('bf'))", u.Password))
	return result.Error
}

// GetUser returns a UserID from the database.
func GetUser(u models.User) (string, error) {
	var id struct{ ID string }
	err := DB.Table("users").Where("username = ?", u.Username).Where(
		"(password_hash = crypt(?, password_hash)) = 't'", u.Password).Select("id").First(&id).Error
	if err != nil {
		return "", err
	}
	log.Println(id)
	return id.ID, nil
}
