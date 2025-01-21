-- migrations/001_create_students_table.up.sql
CREATE TABLE students (
  id BIGSERIAL PRIMARY KEY,
  name text NOT NULL,
  s_class text,
  school text,
  order_day SMALLINT,
  order_time TIME,
  order_cost SMALLINT
);
