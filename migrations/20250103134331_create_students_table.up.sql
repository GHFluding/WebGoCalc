-- migrations/001_create_students_table.up.sql
CREATE TABLE students (
  id BIGSERIAL PRIMARY KEY,
  name text NOT NULL UNIQUE,
  clas text,
  scool text,
  order_day SMALLINT,
  order_time TIME,
  order_cost SMALLINT
);