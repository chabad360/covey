package models

type User struct {
	ID       uint8  `json:"id" gorm:"<-:create;primarykey"`
	Username string `json:"username" gorm:"<-:create;notnull;unique"`
	Password string `json:"password" gorm:"->:false;notnull"`
}
