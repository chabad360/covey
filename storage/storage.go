package storage

import (
	"log"

	"database/sql"

	// Needed for postgres
	_ "github.com/lib/pq"
)

var (
	db *sql.DB
)

// Init initializes the database connection.
func Init() {
	var err error
	db, err = sql.Open("postgres", "user=postgres host=127.0.0.1 port=5432 dbname=covey sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

// GetDB returns the db.
func GetDB() *sql.DB { return db }

// GetItem returns the JSON representation of an item in the database.
func GetItem(table, id string) ([]byte, error) {
	var j []byte
	if err := db.QueryRow("SELECT to_jsonb("+table+") FROM "+table+" WHERE id = $1 OR id_short = $1 OR name = $1;",
		id).Scan(&j); err != nil {
		return nil, err
	}
	return j, nil
}
