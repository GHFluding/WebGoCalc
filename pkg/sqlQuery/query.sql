-- name: GetStudentById :one
SELECT * FROM students 
WHERE id = $1 
LIMIT 1;

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
WHERE order_day = CASE 
    WHEN $2 = 0 THEN 7  
END;

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

-- name: AddEventsForToday :exec
INSERT INTO calendar (student_id, event_date, order_time, order_cost)
SELECT 
    id AS student_id, 
    CURRENT_DATE AS event_date, 
    order_time, 
    order_cost
FROM students
WHERE order_day = CASE 
    WHEN EXTRACT(DOW FROM CURRENT_DATE) = 0 THEN 7 
    ELSE EXTRACT(DOW FROM CURRENT_DATE)
END;
