package safe

// Task is a safe version of models.Task.
type Task struct {
	Plugin  string            `json:"plugin" gorm:"<-:create;notnull"`
	Node    string            `json:"node" gorm:"<-:create;notnull"`
	Details map[string]string `json:"details" gorm:"<-:create;"`
}
