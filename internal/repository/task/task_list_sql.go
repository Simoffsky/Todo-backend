package task

import (
	"context"
	"todo/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresTaskListRepository struct {
	db *pgxpool.Pool
}

func NewPostgresTaskListRepository(db *pgxpool.Pool) *PostgresTaskListRepository {
	return &PostgresTaskListRepository{db: db}
}

func (r *PostgresTaskListRepository) CreateTaskList(list models.TaskList) (int, error) {
	var id int
	err := r.db.QueryRow(context.Background(), `
        INSERT INTO task_lists (owner)
        VALUES ($1)
        RETURNING id
    `, list.Owner).Scan(&id)
	if err != nil {
		return 0, err
	}

	for _, task := range list.Tasks {
		_, err := r.db.Exec(context.Background(), `
            INSERT INTO tasks (title, owner, body, is_done, task_list_id)
            VALUES ($1, $2, $3, $4, $5)
        `, task.Title, task.Owner, task.Body, task.IsDone, id)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

func (r *PostgresTaskListRepository) GetTaskList(id int) (models.TaskList, error) {
	var list models.TaskList
	err := r.db.QueryRow(context.Background(), `
        SELECT id, owner
        FROM task_lists
        WHERE id = $1
    `, id).Scan(&list.ID, &list.Owner)
	if err != nil {
		return models.TaskList{}, models.ErrTaskListNotFound
	}

	rows, err := r.db.Query(context.Background(), `
        SELECT id, title, owner, body, is_done
        FROM tasks
        WHERE task_list_id = $1
    `, id)
	if err != nil {
		return models.TaskList{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Owner, &task.Body, &task.IsDone); err != nil {
			return models.TaskList{}, err
		}
		list.Tasks = append(list.Tasks, task)
	}

	return list, nil
}

func (r *PostgresTaskListRepository) DeleteTaskList(id int) error {
	_, err := r.db.Exec(context.Background(), `
        DELETE FROM task_lists
        WHERE id = $1
    `, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresTaskListRepository) UpdateTaskList(updatedList models.TaskList) error {
	_, err := r.db.Exec(context.Background(), `
        UPDATE task_lists
        SET owner = $2
        WHERE id = $1
    `, updatedList.ID, updatedList.Owner)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(context.Background(), `
        DELETE FROM tasks
        WHERE task_list_id = $1
    `, updatedList.ID)
	if err != nil {
		return err
	}

	for _, task := range updatedList.Tasks {
		_, err := r.db.Exec(context.Background(), `
            INSERT INTO tasks (title, owner, body, is_done, task_list_id)
            VALUES ($1, $2, $3, $4, $5)
        `, task.Title, task.Owner, task.Body, task.IsDone, updatedList.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
