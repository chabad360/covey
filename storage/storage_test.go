package storage

import (
	"github.com/chabad360/covey/config"
	"gorm.io/gorm"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/chabad360/covey/test"
)

func TestInit(t *testing.T) {
	Init()
	if reflect.TypeOf(DB) != reflect.TypeOf(&gorm.DB{}) {
		t.Errorf("Init() = %v, want %v", reflect.TypeOf(gorm.DB{}), reflect.TypeOf(DB))
	}
}

func TestMain(m *testing.M) {
	pool, resource, pdb, err := test.Boilerplate()
	DB = pdb
	if err != nil {
		log.Fatalf("Could not setup DB connection: %s", err)
	}

	config.Config.DB.Username = "postgres"
	config.Config.DB.Password = "secret"
	config.Config.DB.Host = "localhost"
	config.Config.DB.Port = resource.GetPort("5432/tcp")
	config.Config.DB.Database = "covey"
	if err != nil {
		log.Fatalf("Could not setup config")
	}

	DB.Create(&task)
	DB.Create(&j2)

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
