package storage

import (
	"errors"
	"github.com/chabad360/covey/models"
	"gorm.io/gorm"
)

// AddNode adds a node to the database.
func AddNode(node *models.Node) error {
	result := DB.Create(node)
	return result.Error
}

// GetNodeIDorName returns the full ID or name for the given node.
func GetNodeIDorName(id string, field string) (string, bool) {
	var ID string
	result := DB.Table("nodes").Where("id = ?", id).Or("id_short = ?", id).Or("name = ?", id).Select(field).Scan(&ID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", false
	}

	return ID, true
}

// GetNode returns the JSON of a node and its keys separately.
func GetNode(id string) (*models.Node, bool) {
	var n models.Node
	result := DB.Where("id = ?", id).Or("id_short = ?", id).Or("name = ?", id).First(&n)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, false
	}

	return &n, true
}

func DeleteNode(node *models.Node) error {
	result := DB.Delete(node)
	return result.Error
}
