package storage

import (
	"fmt"
	"github.com/chabad360/covey/config"
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
	DB, err = gorm.Open(postgres.Open(fmt.Sprintf("user=%s password=%s host=%s port=%v dbname=%s", config.Config.DB.Username, config.Config.DB.Password, config.Config.DB.Host, config.Config.DB.Port, config.Config.DB.Database)), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Silent),
	})
	return err
}
