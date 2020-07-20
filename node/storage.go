package node

import (
	"errors"
	"github.com/chabad360/covey/models"
	"gorm.io/gorm"

	"github.com/chabad360/covey/storage"
)

var db *gorm.DB

// AddNode adds a node to the database.
func addNode(node *models.Node) error {
	refreshDB()

	result := db.Create(node)
	return result.Error
}

// GetNodeIDorName returns the full ID or name for the given node.
func GetNodeIDorName(id string, field string) (string, bool) {
	refreshDB()

	var ID string
	result := db.Table("nodes").Where("id = ?", id).Or("id_short = ?", id).Or("name = ?", id).Select(field).Scan(&ID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", false
	}

	return ID, true
}

// GetNode returns the JSON of a node and its keys separately.
func GetNode(id string) (*models.Node, bool) {
	refreshDB()

	var n models.Node
	result := db.Where("id = ?", id).Or("id_short = ?", id).Or("name = ?", id).First(&n)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, false
	}

	return &n, true
}

func deleteNode(node *models.Node) error {
	refreshDB()

	result := db.Delete(node)
	return result.Error
}

func refreshDB() {
	if db == nil {
		db = storage.DB
	}
}
