package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
)

func TestGetItem(t *testing.T) {
	var tests = []struct {
		id   string
		want string
	}{
		{"1", `{"id": "1", "name": "1", "plugin": "1", "details": 1, "id_short": "1"}`},
		{"2", ""},
	}

	db.Exec("INSERT INTO nodes(id, id_short, name, plugin, details) VALUES($1, $2, $3, $4, $5);",
		"1", "1", "1", "1", "1")

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.id)
		t.Run(testname, func(t *testing.T) {
			if got, _ := GetItem("nodes", tt.id); string(got) != tt.want {
				t.Errorf("GetItem() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestGetDB(t *testing.T) {
	db := GetDB()
	if reflect.TypeOf(db) != reflect.TypeOf(&sql.DB{}) {
		t.Errorf("GetDB() = %v, want %v", reflect.TypeOf(sql.DB{}), reflect.TypeOf(db))
	}
}

func TestInit(t *testing.T) {
	Init()
	if reflect.TypeOf(db) != reflect.TypeOf(&sql.DB{}) {
		t.Errorf("GetDB() = %v, want %v", reflect.TypeOf(sql.DB{}), reflect.TypeOf(db))
	}
}

func TestMain(m *testing.M) {
	var err error
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("postgres", "9.6", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=covey"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", fmt.Sprintf("postgres://postgres:secret@localhost:%s/%s?sslmode=disable",
			resource.GetPort("5432/tcp"), "covey"))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	db.Exec(`
	CREATE TABLE nodes (
		id TEXT PRIMARY KEY NOT NULL,
		id_short TEXT UNIQUE NOT NULL,
		name TEXT UNIQUE NOT NULL,
		plugin TEXT NOT NULL,
		details JSONB NOT NULL
	);`)

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
