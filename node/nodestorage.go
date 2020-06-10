package node

import (
	"context"

	"github.com/chabad360/covey/node/types"
	"github.com/chabad360/covey/storage"
)

// AddNode adds a node to the database.
func addNode(node types.INode) error {
	db := storage.GetPool()
	_, err := db.Exec(context.Background(), "INSERT INTO nodes(id, id_short, name, plugin, details) VALUES($1, $2, $3, $4, $5);",
		node.GetID(), node.GetIDShort(), node.GetName(), node.GetPlugin(), node.GetDetails())
	return err
}
