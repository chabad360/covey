package test

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest"
)

//revive:disable:function-result-limit

// Boilerplate creates a new dockertest connection.
func Boilerplate() (*dockertest.Pool, *dockertest.Resource, *pgxpool.Pool, error) {
	var err error
	var db *pgxpool.Pool
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not connect to docker: %s", err)
	}

	resource, err := pool.Run("postgres", "9.6", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=covey"})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		var err error
		db, err = pgxpool.Connect(context.Background(),
			fmt.Sprintf("postgres://postgres:secret@localhost:%s/%s?sslmode=disable",
				resource.GetPort("5432/tcp"), "covey"))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, nil, nil, fmt.Errorf("could not connect to docker: %s", err)
	}

	db.Exec(context.Background(), `
	CREATE EXTENSION pgcrypto;

	CREATE TABLE nodes (
		id TEXT PRIMARY KEY NOT NULL,
		id_short TEXT UNIQUE NOT NULL,
		name TEXT UNIQUE NOT NULL,
		plugin TEXT NOT NULL,
		details JSONB NOT NULL
	);
	
	CREATE TABLE tasks (
		id TEXT PRIMARY KEY NOT NULL,
		id_short TEXT UNIQUE NOT NULL,
		plugin TEXT NOT NULL,
		state INT NOT NULL,
		node TEXT NOT NULL,
		time TEXT,
		log JSONB,
		details JSONB NOT NULL
	);
	
	CREATE TABLE jobs (
		id TEXT PRIMARY KEY NOT NULL,
		id_short TEXT UNIQUE NOT NULL,
		name TEXT UNIQUE NOT NULL,
		cron TEXT,
		nodes JSONB NOT NULL,
		tasks JSONB NOT NULL,
		task_history JSONB
	);
	
	CREATE TABLE users (
		id SERIAL PRIMARY KEY NOT NULL,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL
	
	);`)

	return pool, resource, db, nil
}
