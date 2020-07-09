package job

import (
	"errors"
	"github.com/chabad360/covey/models"
	"gorm.io/gorm"

	"github.com/chabad360/covey/storage"
)

var db *gorm.DB

// AddJob adds a Job to the database.
func AddJob(j models.Job) error {
	refreshDB()

	result := db.Create(j)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	return nil
}

// GetJob checks if a job with the identifier exists and returns it.
func GetJob(id string) (*models.Job, bool) {
	refreshDB()

	var j models.Job
	result := db.Where("id = ?", id).Or("id_short = ?", id).Or("name = ?", id).First(&j)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, false
	}

	return &j, true
}

// UpdateJob updates a Job in the database.
func UpdateJob(j models.Job) error {
	refreshDB()
	result := db.Save(j)
	return result.Error
}

//// GetJobWithFullHistory returns a job with the tasks substituted for their IDs.
//// Query designed with the help of https://stackoverflow.com/questions/47275606
//func GetJobWithFullHistory(id string) ([]byte, error) {
//	refreshDB()
//	var b []byte
//	if err := db.QueryRow(context.Background(), `SELECT jsonb_build_object('id', j.id, 'name', j.name, 'cron', j.cron,
//			'nodes', j.nodes, 'tasks', j.tasks, 'task_history', j1.task_history)
//		FROM   jobs j
//			LEFT   JOIN LATERAL (
//			SELECT jsonb_agg(to_jsonb(t) - 'id_short') AS task_history
//			FROM   jsonb_array_elements_text(j.task_history) AS p(id)
//			LEFT   JOIN tasks t ON t.id = p.id
//			GROUP  BY j.id
//		) j1 ON j.task_history <> '[]'
//		WHERE id = $1 OR id_short = $1 OR name = $1;`, id).Scan(&b); err != nil {
//		return nil, err
//	}
//	return b, nil
//}

func refreshDB() {
	if db == nil {
		db = storage.DB
	}
}
