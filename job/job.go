package job

import (
	"github.com/chabad360/covey/models"
	"github.com/chabad360/covey/task"
	"log"

	json "github.com/json-iterator/go"
	"github.com/robfig/cron/v3"
)

var (
	cronTab = cron.New()
)

// Init loads up the the jobs and starts the cronTab.
func Init() {
	refreshDB()
	q, err := db.Table("jobs").Where("cron != ''").Select("id", "cron").Rows()
	if err != nil {
		log.Panic(err)
	}
	defer q.Close()

	for q.Next() {
		var id string
		var c string
		if err = q.Scan(&id, &c); err != nil {
			log.Panic(err)
		}
		if err = addCron(id, c); err != nil {
			log.Panic(err)
		}
	}

	cronTab.Start()
}

func addCron(id string, cron string) error {
	_, err := cronTab.AddFunc(cron, func() {
		if j, ok := getJob(id); ok {
			if _, err := run(j); err != nil {
				log.Panic(err)
			}
		}
	})
	if err != nil {
		return err
	}

	return nil
}

// Run runs each task in succession on the specified nodes (concurrently).
// There seems to be a bug where the tasks can occasionally be sent in the wrong order.
func run(j *models.Job) ([]string, error) {
	var th []string
	for _, t := range j.Tasks {
		for _, node := range j.Nodes {
			t.Node = node
			x, err := json.Marshal(t)
			if err != nil {
				log.Panic(err)
			}

			r, err := task.NewTask(x)
			if err != nil {
				log.Panic(err)
			}

			th = append(th, r.ID)
		}
	}
	j.TaskHistory = append(j.TaskHistory, th...)
	return th, updateJob(j)
}
