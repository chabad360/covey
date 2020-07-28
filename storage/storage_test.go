package storage

import (
	"encoding/hex"
	"github.com/chabad360/covey/config"
	"gorm.io/gorm"
	"log"
	"os"
	"reflect"
	"strconv"
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
	config.Config.DB.Port, err = strconv.Atoi(resource.GetPort("5432/tcp"))
	config.Config.DB.Database = "covey"
	if err != nil {
		log.Fatalf("Could not setup config")
	}

	n.HostKey, _ = hex.DecodeString("0000001365636473612d736861322d6e69737470323536000000086e6973747032353600000041044032b5eed25ed08ec4361d9f7e6a7e27f725d563bc033f777fe2b12bdd61c86c160476c6d080b1361ea4ab9e89ec104051762ecb0a4595f53a16a06c959a0704")

	DB.Create(task)
	//DB.Create(j)
	DB.Create(j2)

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
