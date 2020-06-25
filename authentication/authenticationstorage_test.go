package authentication

import (
	"context"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/chabad360/covey/storage"
	"github.com/chabad360/covey/test"
)

var (
	u = &credentials{
		Username: "user",
		Password: "password",
	}

	uu = &credentials{
		Username: "user",
		Password: "pass",
	}

	u2 = &credentials{
		Username: "user2",
		Password: "password",
	}
)

func TestAddUser(t *testing.T) {
	var tests = []struct {
		id   string
		user credentials
		want string
	}{
		{"1", *u, "1"},
		{"2", *u2, "2"},
	}

	testError := AddUser(*u2)

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			var got int
			if db.QueryRow(context.Background(), `SELECT id FROM users 
			WHERE username = $1 AND (password_hash = crypt($2, password_hash)) = 't';`,
				tt.user.Username, tt.user.Password).Scan(&got); strconv.Itoa(got) != tt.want {
				t.Errorf("AddUser() = %v, want %v, error: %v", strconv.Itoa(got), tt.want, testError)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	testError := UpdateUser(*uu, *u)

	var tests = []struct {
		id   string
		user credentials
		want string
	}{
		{"1", *u, "0"},
		{"2", *uu, "1"},
	}

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			var got int
			if db.QueryRow(context.Background(), `SELECT id FROM users 
			WHERE username = $1 AND (password_hash = crypt($2, password_hash)) = 't';`,
				tt.user.Username, tt.user.Password).Scan(&got); strconv.Itoa(got) != tt.want {
				t.Errorf("UpdateUser() = %v, want %v, error: %v", strconv.Itoa(got), tt.want, testError)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	var tests = []struct {
		id   string
		user credentials
		want string
	}{
		{"1", *u, ""},
		{"2", *uu, "1"},
		{"3", *u2, "2"},
	}

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			if got, err := GetUser(tt.user); got != tt.want {
				t.Errorf("GetUser() = %v, want %v, error: %v", got, tt.want, err)
			}
		})
	}
}

func TestMain(m *testing.M) {
	pool, resource, pdb, err := test.Boilerplate()
	db = pdb
	storage.DB = pdb
	if err != nil {
		log.Fatalf("Could not setup DB connection: %s", err)
	}

	_, err = db.Exec(context.Background(), `INSERT INTO users(username, password_hash) 
		VALUES($1, crypt($2, gen_salt('bf')));`,
		u.Username, u.Password)
	if err != nil {
		log.Fatalf("Could not prepare DB: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
