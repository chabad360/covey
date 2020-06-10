package storage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	// DB is necessary mostly for tests.
	DB *pgxpool.Pool
)

// Init initializes the database connection.
func Init() {
	var err error
	// TODO: provide a method for configuration
	DB, err = pgxpool.Connect(context.Background(), "user=postgres host=127.0.0.1 port=5432 dbname=covey")
	if err != nil {
		log.Fatal(err)
	}
}

// GetItem returns the JSON representation of an item in the database.
func GetItem(table, id string) ([]byte, error) {
	var j []byte
	if err := DB.QueryRow(context.Background(), "SELECT to_jsonb("+table+") FROM "+table+
		" WHERE id = $1 OR id_short = $1 OR name = $1;", id).Scan(&j); err != nil {
		return nil, err
	}
	return j, nil
}
