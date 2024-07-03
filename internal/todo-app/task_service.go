package todoapp

import (
	"todo/internal/models"
	repository "todo/internal/repository/task"
)

type TaskService interface {
	CreateTask(task models.Task) (int, error)
	GetTask(id int) (models.Task, error)
	DeleteTask(id int) error
	UpdateTask(task models.Task) error
}

type TaskServiceDefault struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskServiceDefault {
	return &TaskServiceDefault{
		repo: repo,
	}
}

func (s *TaskServiceDefault) CreateTask(task models.Task) (int, error) {
	return s.repo.CreateTask(task)
}

func (s *TaskServiceDefault) GetTask(id int) (models.Task, error) {
	return s.repo.GetTask(id)
}

func (s *TaskServiceDefault) DeleteTask(id int) error {
	return s.repo.DeleteTask(id)
}

func (s *TaskServiceDefault) UpdateTask(task models.Task) error {
	return s.repo.UpdateTask(task)
}
