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

-- name: GetEventsByDate :many
SELECT 
    c.id AS calendar_id,
    s.id AS student_id,
    s.name AS student_name,
    c.event_date,
    c.order_time,
    c.order_cost,
    c.order_check
FROM calendar c
JOIN students s ON c.student_id = s.id
WHERE c.event_date = $1
ORDER BY c.order_time;

-- name: AddEventsForDay :exec
INSERT INTO calendar (student_id, event_date, order_time, order_cost)
SELECT 
    id AS student_id, 
    $1::DATE AS event_date, 
    order_time, 
    order_cost
FROM students
WHERE order_day = $2;

-- name: DeleteEventsByDate :exec
DELETE FROM calendar
WHERE event_date = $1;

-- name: DeleteEventsByStudent :exec
DELETE FROM calendar
WHERE student_id = $1;

-- name: MarkEventAsChecked :exec
UPDATE calendar
SET order_check = TRUE
WHERE student_id = $1 AND event_date = $2;


--auto-add to date, but i'm not sure if it work's
-- name: AddEventsForToday :exec
INSERT INTO calendar (student_id, event_date, order_time, order_cost)
SELECT 
    id AS student_id, 
    CURRENT_DATE AS event_date, 
    order_time, 
    order_cost
FROM students
WHERE order_day = EXTRACT(DOW FROM CURRENT_DATE);
