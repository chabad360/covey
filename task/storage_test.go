package task

import (
	"github.com/chabad360/covey/models"
	"github.com/chabad360/covey/storage"
	"github.com/chabad360/covey/test"
	"log"
	"os"
	"reflect"
	"testing"
)

var task = &models.Task{
	ID:       "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
	State:    models.StateRunning,
	Plugin:   "test",
	Details:  map[string]string{"test": "test"},
	ExitCode: 0,
}

func TestAddTask(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		id   string
		want models.Task
	}{
		{"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", *task},
		{"3", models.Task{}},
	}
	//revive:enable:line-length-limit

	testError := addTask(task)

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			var got models.Task
			if db.Where("id = ?", tt.id).First(&got); reflect.DeepEqual(got, tt.want) {
				t.Errorf("addTask() = %v, want %v, error: %v", got, tt.want, testError)
			}
		})
	}
}

func TestSaveTask(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		id   string
		want models.Task
	}{
		{"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", *task},
		{"3", models.Task{}},
	}
	//revive:enable:line-length-limit

	tu := &TaskInfo{[]string{"hello", "world"}, 0, "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e"}
	saveTask(tu)

	for _, tt := range tests {
		testName := tt.id
		t.Run(testName, func(t *testing.T) {
			var got models.Task
			if result := db.Where("id = ?", tt.id).First(&got); reflect.DeepEqual(got, tt.want) {
				t.Errorf("saveTask() = %v, want %v, error: %v", got, tt.want, result.Error)
			}
		})
	}
}

func TestGetTask(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		id   string
		want models.Task
	}{
		{"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", *task},
		{"3", models.Task{}},
	}
	//revive:enable:line-length-limit

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			if got, err := getTask(tt.id); reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTask() = %v, want %v, error: %v", got, tt.want, err)
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

	err = db.AutoMigrate(&models.Task{})
	if err != nil {
		log.Fatalf("error preping the database: %s", err)
	}

	db.Create(task)

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
