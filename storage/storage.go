package storage

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	// DB provides the gorm DB interface.
	DB *gorm.DB
)

// Init initializes the database connection.
func Init() error {
	var err error
	// TODO: provide a method for configuration
	DB, err = gorm.Open(postgres.Open("user=postgres host=127.0.0.1 port=5432 dbname=covey"), &gorm.Config{})
	return err
}
