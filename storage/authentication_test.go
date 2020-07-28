package storage

import (
	"github.com/chabad360/covey/models"
	"testing"
)

var (
	u = &models.User{
		Username: "user",
		Password: "password",
	}

	uu = &models.User{
		Username: "user",
		Password: "pass",
	}

	u2 = &models.User{
		Username: "user2",
		Password: "password",
	}
)

func TestAddUser(t *testing.T) {
	var tests = []struct {
		id   string
		user models.User
		want string
	}{
		{"1", *u, "1"},
		{"2", *u2, "2"},
	}

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			testError := AddUser(tt.user)

			var got struct{ ID string }
			if err := DB.Table("users").Where("username = ?", u.Username).Where(
				"(password_hash = crypt(?, password_hash)) = 't'", u.Password).
				Select("id").First(&got).Error; got.ID != tt.want && err != nil {
				t.Errorf("AddUser() = %v, want %v, error: %v", got.ID, tt.want, testError)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	var tests = []struct {
		id   string
		user models.User
		want string
	}{
		{"1", *u, "0"},
		{"2", *uu, "1"},
	}

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			testError := UpdateUser(*uu, *u)

			var got struct{ ID string }
			if err := DB.Table("users").Where("username = ?", u.Username).Where(
				"(password_hash = crypt(?, password_hash)) = 't'", u.Password).
				Select("id").First(&got).Error; got.ID != tt.want && err != nil {
				t.Errorf("UpdateUser() = %v, want %v, error: %v", got.ID, tt.want, testError)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	var tests = []struct {
		id   string
		user *models.User
		want string
	}{
		{"1", u, ""},
		{"2", uu, "1"},
		{"3", u2, "2"},
	}

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			if got, err := GetUser(*tt.user); got != tt.want {
				t.Errorf("GetUser() = %v, want %v, error: %v", got, tt.want, err)
			}
		})
	}
}

//func TestMain(m *testing.M) {
//	pool, resource, pdb, err := test.Boilerplate()
//	DB = pdb
//	if err != nil {
//		log.Fatalf("Could not setup DB connection: %s", err)
//	}
//
//	err = DB.AutoMigrate(&models.User{})
//	if err != nil {
//		log.Fatalf("error preping the database: %s", err)
//	}
//
//	result := DB.Exec(`INSERT INTO users(username, password_hash)
//		VALUES(?, crypt(?, gen_salt('bf')));`,
//		u.Username, u.Password)
//	if result.Error != nil {
//		log.Fatalf("Could not prepare DB: %s", result.Error)
//	}
//
//	code := m.Run()
//
//	// You can't defer this because os.Exit doesn't care for defer
//	if err := pool.Purge(resource); err != nil {
//		log.Fatalf("Could not purge resource: %s", err)
//	}
//
//	os.Exit(code)
//}
