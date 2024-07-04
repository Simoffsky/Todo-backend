package task

import (
	"context"
	"todo/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresTaskRepository struct {
	db *pgxpool.Pool
}

func NewPostgresTaskRepository(db *pgxpool.Pool) *PostgresTaskRepository {
	return &PostgresTaskRepository{db: db}
}

func (r *PostgresTaskRepository) CreateTask(task models.Task) (int, error) {
	var id int
	err := r.db.QueryRow(context.Background(), `
        INSERT INTO tasks (title, owner, body, is_done)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `, task.Title, task.Owner, task.Body, task.IsDone).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PostgresTaskRepository) GetTask(id int) (models.Task, error) {
	var task models.Task
	err := r.db.QueryRow(context.Background(), `
        SELECT id, title, owner, body, is_done, task_list_id
        FROM tasks
        WHERE id = $1
    `, id).Scan(&task.ID, &task.Title, &task.Owner, &task.Body, &task.IsDone, &task.TaskListID)
	if err != nil {
		return models.Task{}, models.ErrTaskNotFound
	}

	return task, nil
}

func (r *PostgresTaskRepository) DeleteTask(id int) error {
	cmdTag, err := r.db.Exec(context.Background(), `
        DELETE FROM tasks
        WHERE id = $1
    `, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return models.ErrTaskNotFound
	}

	return nil
}

func (r *PostgresTaskRepository) UpdateTask(updatedTask models.Task) error {
	cmdTag, err := r.db.Exec(context.Background(), `
        UPDATE tasks
        SET title = $2, owner = $3, body = $4, is_done = $5
        WHERE id = $1
    `, updatedTask.ID, updatedTask.Title, updatedTask.Owner, updatedTask.Body, updatedTask.IsDone)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return models.ErrTaskNotFound
	}

	return nil
}
