package todoapp

import (
	"encoding/json"
	"fmt"
	"todo/internal/models"
	repository "todo/internal/repository/task"

	"github.com/go-redis/redis"
)

type TaskService interface {
	CreateTask(task models.Task) (int, error)
	GetTask(id int) (models.Task, error)
	DeleteTask(id int) error
	UpdateTask(task models.Task) error

	CreateTaskList(list models.TaskList) (int, error)
	GetTaskList(id int) (models.TaskList, error)
	DeleteTaskList(id int) error
	UpdateTaskList(list models.TaskList) error
}

type TaskServiceDefault struct {
	taskRepo     repository.TaskRepository
	taskListRepo repository.TaskListRepository
	redisClient  *redis.Client
}

func NewTaskService(repo repository.TaskRepository, listRepo repository.TaskListRepository, redisClient *redis.Client) *TaskServiceDefault {
	return &TaskServiceDefault{
		taskRepo:     repo,
		taskListRepo: listRepo,
		redisClient:  redisClient,
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
	id, err := s.taskListRepo.CreateTaskList(list)
	if err != nil {
		return 0, err
	}
	list.ID = id

	jsonList, _ := json.Marshal(list)
	s.redisClient.Set(fmt.Sprintf("taskList:%d", id), jsonList, 0)
	return id, nil
}

func (s *TaskServiceDefault) UpdateTaskList(list models.TaskList) error {
	err := s.taskListRepo.UpdateTaskList(list)
	if err != nil {
		return err
	}

	jsonList, _ := json.Marshal(list)
	s.redisClient.Set(fmt.Sprintf("taskList:%d", list.ID), jsonList, 0)
	return nil
}

func (s *TaskServiceDefault) GetTaskList(id int) (models.TaskList, error) {
	val, err := s.redisClient.Get(fmt.Sprintf("taskList:%d", id)).Result()
	if err == nil {
		var list models.TaskList
		err := json.Unmarshal([]byte(val), &list)
		if err != nil {
			return models.TaskList{}, err
		}
		return list, nil
	}
	return s.taskListRepo.GetTaskList(id)
}

func (s *TaskServiceDefault) DeleteTaskList(id int) error {
	err := s.taskListRepo.DeleteTaskList(id)
	if err != nil {
		return err
	}
	s.redisClient.Del(fmt.Sprintf("taskList:%d", id))
	return nil
}
