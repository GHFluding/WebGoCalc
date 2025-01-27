-- name: GetStudentById :one
SELECT * FROM students 
WHERE id = $1 
LIMIT 1;

-- name: GetStudentsById :many
SELECT * FROM students 
WHERE order_day = $1;

-- name: ListStudents :many
SELECT * FROM students
ORDER BY id; 

-- name: CreateStudent :one
INSERT INTO students (
  name, s_class, school, order_day, order_time, order_cost
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateStudentById :exec
UPDATE students
SET 
  name = $2,
  s_class = $3,
  school = $4,
  order_day = $5,
  order_time = $6,
  order_cost = $7
WHERE id = $1;

-- name: DeleteStudentById :exec
DELETE FROM students 
WHERE id = $1;

-- name: AddEventsForDay :exec
INSERT INTO student_events (student_id, event_date, order_time, order_cost)
VALUES ($1, $2::DATE, $3, $4);


-- name: DeleteEventsByDate :exec
DELETE FROM student_events
WHERE event_date = $1;

-- name: GetEventsByDate :many
SELECT 
    c.id AS student_events_id,
    s.id AS student_id,
    s.name AS student_name,
    c.event_date,
    c.order_time,
    c.order_cost,
    c.order_check
FROM student_events c
JOIN students s ON c.student_id = s.id
WHERE c.event_date = $1
ORDER BY c.order_time;

-- name: DeleteEventsByStudent :exec
DELETE FROM student_events
WHERE id = $1;

-- name: MarkEventAsChecked :exec
UPDATE student_events
SET order_check = TRUE
WHERE student_id = $1 AND event_date = $2;


