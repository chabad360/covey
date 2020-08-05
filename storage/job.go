package storage

import (
	"github.com/chabad360/covey/models"
)

// AddJob adds a Job to the database.
func AddJob(j *models.Job) error {
	return DB.Create(j).Error
}

// GetJob checks if a job with the identifier exists and returns it.
func GetJob(id string) (*models.Job, bool) {
	var j models.Job
	if DB.
		Where("id = ?", id).
		Or("id_short = ?", id).
		Or("name = ?", id).
		First(&j).
		Error != nil {
		return nil, false
	}

	return &j, true
}

// UpdateJob updates a Job in the database.
func UpdateJob(j *models.Job) error {
	return DB.Save(j).Error
}

// DeleteJob deletes a Job in the database.
func DeleteJob(j *models.Job) error {
	return DB.Delete(j).Error
}

// GetJobWithFullHistory returns a job with the tasks substituted for their IDs.
// Query designed with the help of https://stackoverflow.com/questions/47275606
func GetJobWithFullHistory(id string) (*models.JobWithTasks, bool) {
	var b models.JobWithTasks
	result := DB.Raw(`SELECT j.id, j.name, j.cron, j.nodes, j.tasks, j1.task_history
		FROM   jobs j
			LEFT   JOIN LATERAL (
			SELECT jsonb_agg(to_jsonb(t) - 'details' - 'log') AS task_history
			FROM   jsonb_array_elements_text(convert_from(j.task_history, 'UTF-8')::jsonb) AS p(id)
			LEFT   JOIN tasks t ON t.id = p.id
			GROUP  BY j.id
		) j1 ON j.task_history <> '[]'
		WHERE id = ? OR id_short = ? OR name = ?;`, id, id, id).First(&b)
	if result.Error != nil {
		return nil, false
	}

	return &b, true
}
