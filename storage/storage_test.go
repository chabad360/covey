package storage

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/chabad360/covey/test"
	"github.com/jackc/pgx/v4/pgxpool"
)

func TestInit(t *testing.T) {
	Init()
	if reflect.TypeOf(DB) != reflect.TypeOf(&pgxpool.Pool{}) {
		t.Errorf("Init() = %v, want %v", reflect.TypeOf(pgxpool.Pool{}), reflect.TypeOf(DB))
	}
}

func TestMain(m *testing.M) {
	pool, resource, pdb, err := test.Boilerplate()
	DB = pdb
	if err != nil {
		log.Fatalf("Could not setup DB connection: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
