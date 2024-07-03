package models

type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	IsDone bool   `json:"is_done"`
	Owner  string `json:"owner"` // login of the user who created the task
}

type TaskList struct {
	ID    int
	Tasks []Task
	Owner string // login of the user who owns the list
}
