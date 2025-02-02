-- migrations/001_create_students_table.up.sql
CREATE TABLE students (
  id BIGSERIAL PRIMARY KEY,
  name text NOT NULL,
  s_class text NOT NULL,
  school text NOT NULL,
  order_day SMALLINT NOT NULL,
  order_time TIME NOT NULL,
  order_cost SMALLINT NOT NULL
);