package task

import (
	"github.com/chabad360/covey/storage"
	"github.com/chabad360/covey/task/types"
)

// AddTask adds a task to the database.
func addTask(task types.ITask) error {
	db := storage.GetDB()
	_, err := db.Exec("INSERT INTO tasks(id, id_short, plugin, state, node, time, log, details) VALUES($1, $2, $3, $4, $5, $6, $7, $8);",
		task.GetID(), task.GetIDShort(), task.GetPlugin(), task.GetState(), task.GetNode(), func() string { t, _ := task.GetTime().MarshalText(); return string(t) }(), task.GetLog(), task.GetDetails())
	return err
}

// GetTask returns the JSON representation of a task in the database.
func getTask(id string) (*[]byte, error) {
	db := storage.GetDB()
	var j []byte
	if err := db.QueryRow("SELECT to_jsonb(tasks) FROM tasks WHERE id = $1 OR id_short = $1;", id).Scan(&j); err != nil {
		return nil, err
	}
	return &j, nil
}

// UpdateTask updates a task in the database.
func updateTask(task types.ITask) error {
	db := storage.GetDB()
	_, err := db.Exec("UPDATE tasks SET state = $1, log = $2, details = $3 WHERE id = $4;",
		task.GetState(), task.GetLog(), task.GetDetails(), task.GetID())
	return err
}
