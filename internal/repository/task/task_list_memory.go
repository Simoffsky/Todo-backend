package task

import (
	"sync"
	"todo/internal/models"
)

type InMemoryTaskListRepository struct {
	taskLists map[int]models.TaskList
	mu        sync.RWMutex
	nextID    int
}

func NewInMemoryTaskListRepository() *InMemoryTaskListRepository {
	return &InMemoryTaskListRepository{
		taskLists: make(map[int]models.TaskList),
		nextID:    1,
	}
}

func (r *InMemoryTaskListRepository) CreateTaskList(list models.TaskList) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	list.ID = r.nextID
	r.taskLists[list.ID] = list
	r.nextID++

	return list.ID, nil
}

func (r *InMemoryTaskListRepository) GetTaskList(id int) (models.TaskList, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list, exists := r.taskLists[id]
	if !exists {
		return models.TaskList{}, models.ErrTaskListNotFound
	}

	return list, nil
}

func (r *InMemoryTaskListRepository) DeleteTaskList(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.taskLists[id]
	if !exists {
		return models.ErrTaskListNotFound
	}

	delete(r.taskLists, id)
	return nil
}

func (r *InMemoryTaskListRepository) UpdateTaskList(updatedList models.TaskList) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.taskLists[updatedList.ID]
	if !exists {
		return models.ErrTaskListNotFound
	}

	r.taskLists[updatedList.ID] = updatedList
	return nil
}
