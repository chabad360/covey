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
