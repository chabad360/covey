package job

import (
	"context"
	"log"

	"github.com/chabad360/covey/job/types"
	"github.com/chabad360/covey/storage"
	json "github.com/json-iterator/go"
	"github.com/robfig/cron/v3"
)

var (
	// jobs      = make(map[string]*types.Job)
	// jobsShort = make(map[string]string)
	// jobsName  = make(map[string]string)
	cronTab = cron.New()
)

// Init loads up the the jobs and starts the cronTab.
func Init() {
	db := storage.DB
	q, err := db.Query(context.Background(), "SELECT id, cron FROM jobs WHERE cron != '';")
	if err != nil {
		log.Panic(err)
	}
	defer q.Close()

	for q.Next() {
		var id string
		var cron string
		if err = q.Scan(&id, &cron); err != nil {
			log.Panic(err)
		}
		if err = addCron(id, cron); err != nil {
			log.Panic(err)
		}
	}

	cronTab.Start()
}

// GetJob checks if a job with the identifier exists and returns it.
func GetJob(identifier string) (*types.Job, bool) {
	var t types.Job
	j, err := storage.GetItem("jobs", identifier)
	if err != nil {
		log.Println(err)
		return nil, false
	}
	err = json.Unmarshal(j, &t)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	return &t, true
}

// GetJobWithTasks checks if a job with the identifier exists and returns it along with its tasks.
func GetJobWithTasks(identifier string) (*types.JobWithTasks, bool) {
	var t types.JobWithTasks
	j, err := GetJobWithFullHistory(identifier)
	if err != nil {
		log.Println(err)
		return nil, false
	}
	if err = json.Unmarshal(j, &t); err != nil {
		log.Println(err)
		return nil, false
	}

	return &t, true
}

func addCron(id string, cron string) error {
	_, err := cronTab.AddFunc(cron, func() func() {
		return func() { // This little bundle of joy allows the job to occur despite not being an object.
			j, _ := GetJob(id)
			j.Run()
			if err := UpdateJob(*j); err != nil {
				log.Panic(err)
			}
		}
	}())
	if err != nil {
		return err
	}
	return nil
}
