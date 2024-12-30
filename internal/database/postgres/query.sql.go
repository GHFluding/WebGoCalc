// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createStudent = `-- name: CreateStudent :one
INSERT INTO students (
  name, clas, scool, order_day, order_time, order_cost
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING id, name, clas, scool, order_day, order_time, order_cost
`

type CreateStudentParams struct {
	Name      string
	Clas      pgtype.Text
	Scool     pgtype.Text
	OrderDay  pgtype.Int2
	OrderTime pgtype.Time
	OrderCost pgtype.Int2
}

func (q *Queries) CreateStudent(ctx context.Context, arg CreateStudentParams) (Student, error) {
	row := q.db.QueryRow(ctx, createStudent,
		arg.Name,
		arg.Clas,
		arg.Scool,
		arg.OrderDay,
		arg.OrderTime,
		arg.OrderCost,
	)
	var i Student
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Clas,
		&i.Scool,
		&i.OrderDay,
		&i.OrderTime,
		&i.OrderCost,
	)
	return i, err
}

const deleteStudentByName = `-- name: DeleteStudentByName :exec
DELETE FROM students 
WHERE name = $1
`

func (q *Queries) DeleteStudentByName(ctx context.Context, studentName string) error {
	_, err := q.db.Exec(ctx, deleteStudentByName, studentName)
	return err
}

const getStudentByName = `-- name: GetStudentByName :one
SELECT id, name, clas, scool, order_day, order_time, order_cost FROM students 
WHERE name = $1 
LIMIT 1
`

func (q *Queries) GetStudentByName(ctx context.Context, studentName string) (Student, error) {
	row := q.db.QueryRow(ctx, getStudentByName, studentName)
	var i Student
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Clas,
		&i.Scool,
		&i.OrderDay,
		&i.OrderTime,
		&i.OrderCost,
	)
	return i, err
}

const listStudents = `-- name: ListStudents :many
SELECT id, name, clas, scool, order_day, order_time, order_cost FROM students
ORDER BY name
`

func (q *Queries) ListStudents(ctx context.Context) ([]Student, error) {
	rows, err := q.db.Query(ctx, listStudents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Student
	for rows.Next() {
		var i Student
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Clas,
			&i.Scool,
			&i.OrderDay,
			&i.OrderTime,
			&i.OrderCost,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateStudent = `-- name: UpdateStudent :exec
UPDATE students
SET 
  clas = $2,
  scool = $3,
  order_day = $4,
  order_time = $5,
  order_cost = $6
WHERE id = $1
`

type UpdateStudentParams struct {
	ID        int64
	Clas      pgtype.Text
	Scool     pgtype.Text
	OrderDay  pgtype.Int2
	OrderTime pgtype.Time
	OrderCost pgtype.Int2
}

func (q *Queries) UpdateStudent(ctx context.Context, arg UpdateStudentParams) error {
	_, err := q.db.Exec(ctx, updateStudent,
		arg.ID,
		arg.Clas,
		arg.Scool,
		arg.OrderDay,
		arg.OrderTime,
		arg.OrderCost,
	)
	return err
}
