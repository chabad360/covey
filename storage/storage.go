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
	_, err := db.Exec(context.Background(), "insert into "+table+"(id, id_short, data) values($1, $2, $3);", id, idShort, item)
	return err
}

// GetItem returns an interface of an item in the database.
func GetItem(table string, id string, i interface{}) (interface{}, error) {
	if err := db.QueryRow(context.Background(), "select data from "+table+" where id = $1 or id_short = $1 or data->>'name' = $1;", id).Scan(&i); err != nil {
		return nil, err
	}
	return i, nil
}
