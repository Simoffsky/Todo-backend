package task

import "todo/internal/models"

type TaskRepository interface {
	CreateTask(task models.Task) (int, error)
	GetTask(id int) (models.Task, error)
	DeleteTask(id int) error
}


