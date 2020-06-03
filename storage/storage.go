package storage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	db *pgxpool.Pool
)

// Init initializes the database connection.
func Init() {
	var err error
	db, err = pgxpool.Connect(context.Background(), "user=postgres host=127.0.0.1 port=5432 dbname=covey")
	if err != nil {
		log.Fatal(err)
	}

}

// AddItem adds an item to the database.
func AddItem(table string, id string, idShort string, item interface{}) error {
	_, err := db.Exec(context.Background(), "INSERT INTO "+table+"(id, id_short, data) VALUES($1, $2, $3);", id, idShort, item)
	return err
}

// GetItem returns the JSON representation of an item in the database.
func GetItem(table, id string) (*[]byte, error) {
	var j []byte
	if err := db.QueryRow(context.Background(), "SELECT to_jsonb("+table+") FROM "+table+" WHERE id = $1 OR id_short = $1 OR name = $1;", id).Scan(&j); err != nil {
		return nil, err
	}
	return &j, nil
}

// UpdateItem updates an item in the database.
func UpdateItem(table string, id string, i interface{}) error {
	_, err := db.Exec(context.Background(), "UPDATE "+table+" SET data = $1 WHERE id = $2 OR id_short = $2 OR data->>'name' = $2;", i, id)
	return err
}

// GetJob returns a job with the tasks subsituted for their IDs.
func GetJob(id string, i interface{}) (interface{}, error) { // Query designed with the help of https://stackoverflow.com/questions/47275606
	if err := db.QueryRow(context.Background(), `SELECT COALESCE(j1.data, j.data) 
	FROM   jobs j
	LEFT   JOIN LATERAL (
	   SELECT j.data || jsonb_build_object('task_history', COALESCE(ts.task, '[]')) AS data
	   FROM   jobs
	   CROSS  JOIN LATERAL (
		  SELECT jsonb_agg(to_jsonb(t)) AS task
		  FROM   jsonb_array_elements_text(jobs.data->'task_history') AS p(id)
		  LEFT   JOIN tasks t ON t.id = p.id
		  ) ts
	   ) j1 ON j.data <> '{}'
	WHERE j.id = $1 OR j.id_short = $1 OR j.data->>'name' = $1;`, id).Scan(&i); err != nil {
		return nil, err
	}
	return i, nil
}
