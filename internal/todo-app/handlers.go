package todoapp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"todo/internal/models"
)

func (a *App) handleTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.handleGetTask(w, r)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := w.Write([]byte("Method not allowed"))
		if err != nil {
			a.handleError(w, models.NewError(err, http.StatusInternalServerError))
			return
		}
	}
}

// returns created task as JSON
func (a *App) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		a.handleError(w, models.NewError(err, http.StatusBadRequest))
		return
	}

	id, err := a.taskRepository.CreateTask(task)
	if err != nil {
		a.handleError(w, err)
		return
	}

	task.ID = id

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(task)

	if err != nil {
		a.handleError(w, models.NewError(err, http.StatusInternalServerError))
		return
	}

}

func (a *App) handleGetTask(w http.ResponseWriter, r *http.Request) {
	idPath := r.PathValue("task_id")
	id, err := strconv.Atoi(idPath)
	if err != nil {
		a.handleError(w, models.NewError(err, http.StatusBadRequest))
		return
	}

	task, err := a.taskRepository.GetTask(id)
	if err != nil {
		a.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		a.handleError(w, models.NewError(err, http.StatusInternalServerError))
		return
	}
}


func (a *App) handleError(w http.ResponseWriter, err error) {
	var modelErr models.Error
	if !errors.As(err, &modelErr) {
		a.writeError(w, http.StatusInternalServerError, err)
		return
	}
	a.writeError(w, modelErr.StatusCode, err)
}

func (a *App) writeError(w http.ResponseWriter, statusCode int, err error) {
	a.logger.Error(fmt.Sprintf("HTTP error(%d): %s", statusCode, err.Error()))
	http.Error(w, err.Error(), statusCode)
}
