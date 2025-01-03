-- name: GetStudentByName :one
SELECT * FROM students 
WHERE name = @student_name 
LIMIT 1;


-- name: ListStudents :many
SELECT * FROM students
ORDER BY name;

-- name: CreateStudent :one
INSERT INTO students (
  name, clas, scool, order_day, order_time, order_cost
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateStudent :exec
UPDATE students
SET 
  clas = $2,
  scool = $3,
  order_day = $4,
  order_time = $5,
  order_cost = $6
WHERE id = $1;

-- name: DeleteStudentByName :exec
DELETE FROM students 
WHERE name = @student_name;
