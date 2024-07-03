package task

import (
	"sync"
	"todo/internal/models"
)

type InMemoryTaskRepository struct {
	tasks  map[int]models.Task
	mu     sync.RWMutex
	nextID int
}

func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		tasks:  make(map[int]models.Task),
		nextID: 1,
	}
}

func (r *InMemoryTaskRepository) CreateTask(task models.Task) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task.ID = r.nextID
	r.tasks[task.ID] = task
	r.nextID++

	return task.ID, nil
}

func (r *InMemoryTaskRepository) GetTask(id int) (models.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return models.Task{}, models.ErrTaskNotFound
	}

	return task, nil
}

func (r *InMemoryTaskRepository) DeleteTask(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.tasks[id]
	if !exists {
		return models.ErrTaskNotFound
	}

	delete(r.tasks, id)
	return nil
}
