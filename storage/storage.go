package storage

import (
	"fmt"
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
	DB, err = gorm.Open(postgres.Open(fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", "postgres", "", "127.0.0.1", "5432", "covey")), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Silent),
	})
	return err
}
