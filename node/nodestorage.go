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
		`INSERT INTO nodes(id, id_short, name, private_key, public_key, host_key, username, port, ip) 
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
		node.GetID(), node.GetIDShort(), node.GetName(),
		node.PrivateKey, node.PublicKey, node.HostKey, node.Username, node.Port, node.IP)
	return err
}

// GetNodeID returns the full ID for the given node.
func GetNodeID(identifier string) (string, bool) {
	refreshDB()
	var ID string
	err := db.QueryRow(context.Background(),
		"SELECT id FROM nodes WHERE id = $1 OR id_short = $1 OR name = $1;", identifier).Scan(&ID)
	if err != nil {
		return "", false
	}
	return ID, true
}

func refreshDB() {
	if db == nil {
		db = storage.DB
	}
}
