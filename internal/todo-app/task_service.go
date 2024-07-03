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

	CreateTaskList(list models.TaskList) (int, error)
	GetTaskList(id int) (models.TaskList, error)
	DeleteTaskList(id int) error
}

type TaskServiceDefault struct {
	taskRepo     repository.TaskRepository
	taskListRepo repository.TaskListRepository
}

func NewTaskService(repo repository.TaskRepository, listRepo repository.TaskListRepository) *TaskServiceDefault {
	return &TaskServiceDefault{
		taskRepo:     repo,
		taskListRepo: listRepo,
	}
}

func (s *TaskServiceDefault) CreateTask(task models.Task) (int, error) {
	return s.taskRepo.CreateTask(task)
}

func (s *TaskServiceDefault) GetTask(id int) (models.Task, error) {
	return s.taskRepo.GetTask(id)
}

func (s *TaskServiceDefault) DeleteTask(id int) error {
	return s.taskRepo.DeleteTask(id)
}

func (s *TaskServiceDefault) UpdateTask(task models.Task) error {
	return s.taskRepo.UpdateTask(task)
}

func (s *TaskServiceDefault) CreateTaskList(list models.TaskList) (int, error) {
	return s.taskListRepo.CreateTaskList(list)
}

func (s *TaskServiceDefault) GetTaskList(id int) (models.TaskList, error) {
	return s.taskListRepo.GetTaskList(id)
}

func (s *TaskServiceDefault) DeleteTaskList(id int) error {
	return s.taskListRepo.DeleteTaskList(id)
}

func (s *TaskServiceDefault) UpdateTaskList(list models.TaskList) error {
	return s.taskListRepo.UpdateTaskList(list)
}
