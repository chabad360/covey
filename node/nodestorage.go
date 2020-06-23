package node

import (
	"context"

	"github.com/chabad360/covey/node/types"
	"github.com/chabad360/covey/storage"
	"github.com/jackc/pgx/v4/pgxpool"
)

var db *pgxpool.Pool

// AddNode adds a node to the database.
func addNode(node *types.Node) error {
	refreshDB()
	_, err := db.Exec(context.Background(),
		"INSERT INTO nodes(id, id_short, name, host_key, private_key, public_key, username, port, ip) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9);",
		node.GetID(), node.GetIDShort(), node.GetName(), node.HostKey, node.PrivateKey, node.PublicKey, node.Username, node.Port, node.IP)
	return err
}

func refreshDB() {
	if db == nil {
		db = storage.DB
	}
}
