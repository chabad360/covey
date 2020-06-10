package job

import (
	"context"

	"github.com/chabad360/covey/job/types"
	"github.com/chabad360/covey/storage"
	"github.com/jackc/pgx/v4/pgxpool"
)

var db *pgxpool.Pool

// AddJob adds a Job to the database.
func AddJob(j types.Job) error {
	refreshDB()
	_, err := db.Exec(context.Background(), `INSERT INTO jobs(id, id_short, name, cron, nodes, tasks, task_history) 
		VALUES($1, $2, $3, $4, $5, $6, $7);`,
		j.GetID(), j.GetIDShort(), j.Name, j.Cron, j.Nodes, j.Tasks, j.TaskHistory)
	return err
}

// UpdateJob updates a Job in the database.
func UpdateJob(j types.Job) error {
	refreshDB()
	_, err := db.Exec(context.Background(),
		"UPDATE jobs SET name = $1, cron = $2, nodes = $3, tasks = $4, task_history = $5 WHERE id = $6;",
		j.Name, j.Cron, j.Nodes, j.Tasks, j.TaskHistory, j.GetID())
	return err
}

// GetJobWithFullHistory returns a job with the tasks subsituted for their IDs.
// Query designed with the help of https://stackoverflow.com/questions/47275606
func GetJobWithFullHistory(id string) ([]byte, error) {
	refreshDB()
	var b []byte
	if err := db.QueryRow(context.Background(), `SELECT jsonb_build_object('id', j.id, 'name', j.name, 'cron', j.cron, 
			'nodes', j.nodes, 'tasks', j.tasks, 'task_history', j1.task_history)
		FROM   jobs j
			LEFT   JOIN LATERAL (
			SELECT jsonb_agg(to_jsonb(t) - 'id_short') AS task_history
			FROM   jsonb_array_elements_text(j.task_history) AS p(id)
			LEFT   JOIN tasks t ON t.id = p.id
			GROUP  BY j.id
		) j1 ON j.task_history <> '[]'
		WHERE id = $1 OR id_short = $1 OR name = $1;`, id).Scan(&b); err != nil {
		return nil, err
	}
	return b, nil
}

func refreshDB() {
	if db == nil {
		db = storage.DB
	}
}
