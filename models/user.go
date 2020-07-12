package models

// User contains the info for a user.
type User struct {
	ID       uint8  `json:"id" gorm:"<-:create;primarykey;type:serial"`
	Username string `json:"username" gorm:"<-:create;notnull;unique"`
	Password string `json:"password" gorm:"->:false;notnull;column:password_hash"`
}
