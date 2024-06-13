// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: tasks.sql

package db

import (
	"context"
)

const createTask = `-- name: CreateTask :one
INSERT INTO tasks (title, content, user_id)
VALUES ($1, $2, $3)
RETURNING id, title, content, user_id
`

type CreateTaskParams struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  int32  `json:"user_id"`
}

func (q *Queries) CreateTask(ctx context.Context, arg CreateTaskParams) (Task, error) {
	row := q.db.QueryRowContext(ctx, createTask, arg.Title, arg.Content, arg.UserID)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Content,
		&i.UserID,
	)
	return i, err
}

const deleteTask = `-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1 AND user_id = $2
`

type DeleteTaskParams struct {
	ID     int32 `json:"id"`
	UserID int32 `json:"user_id"`
}

func (q *Queries) DeleteTask(ctx context.Context, arg DeleteTaskParams) error {
	_, err := q.db.ExecContext(ctx, deleteTask, arg.ID, arg.UserID)
	return err
}

const getTaskByUserID = `-- name: GetTaskByUserID :many
SELECT id, title, content, user_id
FROM tasks
WHERE user_id = $1
`

func (q *Queries) GetTaskByUserID(ctx context.Context, userID int32) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, getTaskByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Task{}
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Content,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTask = `-- name: UpdateTask :exec
UPDATE tasks
SET title = $2, content = $3
WHERE id = $1 AND user_id = $4
`

type UpdateTaskParams struct {
	ID      int32  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  int32  `json:"user_id"`
}

func (q *Queries) UpdateTask(ctx context.Context, arg UpdateTaskParams) error {
	_, err := q.db.ExecContext(ctx, updateTask,
		arg.ID,
		arg.Title,
		arg.Content,
		arg.UserID,
	)
	return err
}
