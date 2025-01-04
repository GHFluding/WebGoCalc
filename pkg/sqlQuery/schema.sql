CREATE TABLE students (
  id BIGSERIAL PRIMARY KEY,
  name text NOT NULL UNIQUE,
  s_class text,
  school text,
  order_day SMALLINT,
  order_time TIME,
  order_cost SMALLINT
);
