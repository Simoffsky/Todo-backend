-- +goose Up
BEGIN;
CREATE TABLE task_lists (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    owner VARCHAR(255) REFERENCES users(login) ON DELETE CASCADE,
);

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    owner VARCHAR(255) REFERENCES users(login) ON DELETE CASCADE,
    body TEXT,
    is_done BOOLEAN DEFAULT FALSE,
    task_list_id INT,
    FOREIGN KEY (task_list_id) REFERENCES task_lists(id) ON DELETE CASCADE
);


COMMIT;

-- +goose Down
BEGIN;

DROP TABLE tasks;
DROP TABLE task_lists;

COMMIT;

