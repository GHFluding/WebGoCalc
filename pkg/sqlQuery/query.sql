-- name: GetStudentByName :one
SELECT * FROM students 
WHERE name = $1 
LIMIT 1;

-- name: ListStudents :many
SELECT * FROM students
ORDER BY name;

-- name: CreateStudent :one
INSERT INTO students (
  name, s_class, school, order_day, order_time, order_cost
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateStudentByName :exec
UPDATE students
SET 
  s_class = $2,
  school = $3,
  order_day = $4,
  order_time = $5,
  order_cost = $6
WHERE name = $1;

-- name: DeleteStudentByName :exec
DELETE FROM students 
WHERE name = $1;
