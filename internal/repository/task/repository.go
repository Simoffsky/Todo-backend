package task

import "todo/internal/models"

type TaskRepository interface {
	CreateTask(task models.Task) (int, error)
	GetTask(id int) (models.Task, error)
	DeleteTask(id int) error
	UpdateTask(task models.Task) error
}

type TaskListRepository interface {
	CreateTaskList(list models.TaskList) (int, error)
	GetTaskList(id int) (models.TaskList, error)
	DeleteTaskList(id int) error
	UpdateTaskList(list models.TaskList) error
}
