package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/chabad360/covey/test"
	"github.com/jackc/pgx/v4/pgxpool"
)

func TestGetItem(t *testing.T) {
	var tests = []struct {
		id   string
		want string
	}{
		{"1", `{"id": "1", "name": "1", "plugin": "1", "details": 1, "id_short": "1"}`},
		{"2", ""},
	}

	DB.Exec(context.Background(), "INSERT INTO nodes(id, id_short, name, plugin, details) VALUES($1, $2, $3, $4, $5);",
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
