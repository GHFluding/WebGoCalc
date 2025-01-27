CREATE TABLE students (
  id BIGSERIAL PRIMARY KEY,
  name text NOT NULL,
  s_class text NOT NULL,
  school text NOT NULL,
  order_day SMALLINT NOT NULL,
  order_time TIME NOT NULL,
  order_cost SMALLINT NOT NULL
);

CREATE TABLE student_events (
    id BIGSERIAL PRIMARY KEY,
    student_id BIGINT NOT NULL,
    event_date DATE NOT NULL,
    order_time TIME NOT NULL,
    order_cost SMALLINT NOT NULL,
    order_check BOOLEAN DEFAULT FALSE
);
