package storage

import (
	"context"

	"github.com/chabad360/covey/node/types"
)

// AddNode adds a node to the database.
func AddNode(node types.INode) error {
	_, err := db.Exec(context.Background(), "INSERT INTO nodes(id, id_short, name, plugin, details) VALUES($1, $2, $3, $4, $5);",
		node.GetID(), node.GetIDShort(), node.GetName(), node.GetPlugin(), node.GetDetails())
	return err
}

// GetNode returns the JSON representation of a node in the database.
func GetNode(id string) (*[]byte, error) {
	var j []byte
	if err := db.QueryRow(context.Background(), "SELECT to_jsonb(nodes) FROM nodes WHERE id = $1 OR id_short = $1 OR name = $1;", id).Scan(&j); err != nil {
		return nil, err
	}
	return &j, nil
}
