// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package storage

import (
	"context"
)

const createUser = `-- name: CreateUser :exec
INSERT INTO users (id, username, password) VALUES (?, ?, ?)
`

type CreateUserParams struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.exec(ctx, q.createUserStmt, createUser, arg.ID, arg.Username, arg.Password)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.deleteUserStmt, deleteUser, id)
	return err
}

const getPasswordByID = `-- name: GetPasswordByID :one
SELECT password FROM users WHERE id = ?
`

func (q *Queries) GetPasswordByID(ctx context.Context, id int64) (string, error) {
	row := q.queryRow(ctx, q.getPasswordByIDStmt, getPasswordByID, id)
	var password string
	err := row.Scan(&password)
	return password, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, username FROM users WHERE id = ?
`

type GetUserByIDRow struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

func (q *Queries) GetUserByID(ctx context.Context, id int64) (GetUserByIDRow, error) {
	row := q.queryRow(ctx, q.getUserByIDStmt, getUserByID, id)
	var i GetUserByIDRow
	err := row.Scan(&i.ID, &i.Username)
	return i, err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users SET username = ? WHERE id = ?
`

type UpdateUserParams struct {
	Username string `json:"username"`
	ID       int64  `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.exec(ctx, q.updateUserStmt, updateUser, arg.Username, arg.ID)
	return err
}
