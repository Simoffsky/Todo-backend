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
	case http.MethodDelete:
		a.handleDeleteTask(w, r)
	case http.MethodPut:
		a.handleUpdateTask(w, r)
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

	login, err := getLoginByRequest(r)
	if err != nil {
		a.handleError(w, err)
		return
	}
	task.Owner = login

	id, err := a.taskService.CreateTask(task)
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
	id, err := getIdByRequest(r)
	if err != nil {
		a.handleError(w, err)
		return
	}

	err = a.checkPermission(r, id)
	if err != nil {
		a.handleError(w, err)
		return
	}

	task, err := a.taskService.GetTask(id)
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

func (a *App) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := getIdByRequest(r)
	if err != nil {
		a.handleError(w, err)
		return
	}

	err = a.checkPermission(r, id)
	if err != nil {
		a.handleError(w, err)
		return
	}

	err = a.taskService.DeleteTask(id)
	if err != nil {
		a.handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *App) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := getIdByRequest(r)
	if err != nil {
		a.handleError(w, err)
		return
	}

	err = a.checkPermission(r, id)
	if err != nil {
		a.handleError(w, err)
		return
	}

	var task models.Task
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		a.handleError(w, models.NewError(err, http.StatusBadRequest))
		return
	}

	task.ID = id

	err = a.taskService.UpdateTask(task)
	if err != nil {
		a.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *App) handleRegister(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		a.writeError(w, http.StatusBadRequest, err)
		return
	}

	a.logger.Debug("sending gRPC request(register) with login: " + user.Login)
	if err := a.authService.Register(user); err != nil {
		a.handleError(w, err)
		return
	}
	fmt.Println("Registered successfully")
	w.WriteHeader(http.StatusCreated)
}

func (a *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		a.writeError(w, http.StatusBadRequest, err)
		return
	}

	a.logger.Debug("sending gRPC request(login) with login: " + user.Login)
	token, err := a.authService.Login(user)
	if err != nil {
		a.handleError(w, err)
		return
	}

	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusOK)

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

func (a *App) checkPermission(r *http.Request, taskId int) error {
	login, err := getLoginByRequest(r)
	if err != nil {
		return err
	}

	task, err := a.taskService.GetTask(taskId)
	if err != nil {
		return err
	}

	if task.Owner != login {
		return models.ErrAccessDenied
	}

	return nil
}
func getIdByRequest(r *http.Request) (int, error) {
	idPath := r.PathValue("task_id")
	id, err := strconv.Atoi(idPath)
	if err != nil {
		return 0, models.NewError(err, http.StatusBadRequest)
	}
	return id, nil
}

func getLoginByRequest(r *http.Request) (string, error) {
	login, ok := r.Context().Value(LoginKey("login")).(string)
	if !ok {
		return "", models.NewError(errors.New("login not found in context"), http.StatusInternalServerError)
	}
	return login, nil
}
