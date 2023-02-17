// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: auth.sql

package orm

import (
	"context"
)

const userTestPassword = `-- name: UserTestPassword :one
SELECT user_test_password($1, $2)
`

type UserTestPasswordParams struct {
	Usern string
	Pass  string
}

func (q *Queries) UserTestPassword(ctx context.Context, arg UserTestPasswordParams) (bool, error) {
	row := q.db.QueryRow(ctx, userTestPassword, arg.Usern, arg.Pass)
	var user_test_password bool
	err := row.Scan(&user_test_password)
	return user_test_password, err
}