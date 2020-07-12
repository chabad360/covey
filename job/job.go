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

//// GetJobWithTasks checks if a job with the identifier exists and returns it along with its tasks.
//func GetJobWithTasks(identifier string) (*models.Job, bool) {
//	var t models.Job
//	j, err := GetJobWithFullHistory(identifier)
//	if err != nil {
//		log.Println(err)
//		return nil, false
//	}
//	if err = json.Unmarshal(j, &t); err != nil {
//		log.Println(err)
//		return nil, false
//	}
//
//	return &t, true
//}

func addCron(id string, cron string) error {
	_, err := cronTab.AddFunc(cron, func() func() {
		return func() { // This little bundle of joy allows the job to occur despite not being an object.
			j, _ := GetJob(id)
			Run(j)
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

// Run runs each task in succession on the specified nodes (concurrently).
// There seems to be a bug where the tasks can occasionally be sent in the wrong order
func Run(j *models.Job) {
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

			j.TaskHistory = append(j.TaskHistory, r.ID)
		}
	}
	UpdateJob(*j)
}
