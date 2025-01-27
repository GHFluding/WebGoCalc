-- migrations/001_add_day_table.up.sql
CREATE TABLE calendar (
    id BIGSERIAL PRIMARY KEY,
    student_id BIGINT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    event_date DATE NOT NULL,
    order_time TIME NOT NULL,
    order_cost SMALLINT NOT NULL,
    order_check BOOLEAN DEFAULT FALSE
);