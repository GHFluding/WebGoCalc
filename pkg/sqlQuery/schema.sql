CREATE TABLE students (
  id BIGSERIAL PRIMARY KEY,
  name text NOT NULL,
  s_class text,
  school text,
  order_day SMALLINT,
  order_time TIME,
  order_cost SMALLINT
);

CREATE TABLE calendar (
    id BIGSERIAL PRIMARY KEY,
    student_id BIGINT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    event_date DATE NOT NULL,
    order_time TIME NOT NULL,
    order_cost SMALLINT NOT NULL,
    order_check BOOLEAN DEFAULT FALSE
);
