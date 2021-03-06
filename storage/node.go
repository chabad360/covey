package storage

import (
	"github.com/chabad360/covey/models"
)

// AddNode adds a node to the database.
func AddNode(node *models.Node) error {
	return DB.Create(node).Error
}

// GetNodeIDorName returns the full ID or name for the given node.
func GetNodeIDorName(id string, field string) (string, bool) {
	var ID []string
	if err := DB.Table("nodes").Where("id = ?", id).Or("id_short = ?", id).
		Or("name = ?", id).Pluck(field, &ID).Error; err != nil || len(ID) == 0 {
		return "", false
	}

	return ID[0], true
}

// GetNode returns the JSON of a node and its keys separately.
func GetNode(id string) (*models.Node, bool) {
	var n models.Node
	result := DB.Where("id = ?", id).Or("id_short = ?", id).Or("name = ?", id).First(&n)
	if result.Error != nil {
		return nil, false
	}

	return &n, true
}

// DeleteNode deletes a node from the database.
func DeleteNode(node *models.Node) error {
	return DB.Delete(node).Error
}
