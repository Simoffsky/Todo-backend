package models

type Task struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	IsDone     bool   `json:"is_done"`
	Owner      string `json:"owner"` // login of the user who created the task
	TaskListID int    `json:"task_list_id,omitempty"`
}

type TaskList struct {
	ID    int    `json:"id"`
	Tasks []Task `json:"tasks"`
	Owner string `json:"owner"` // login of the user who owns the list
}
