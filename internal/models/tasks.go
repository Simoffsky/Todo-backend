package models

type Task struct {
	ID     int
	Title  string
	Body   string
	IsDone bool
}

type TaskList struct {
	ID    int
	Tasks []Task
}
