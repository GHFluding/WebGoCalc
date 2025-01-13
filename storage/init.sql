CREATE DATABASE students WITH OWNER postgres TEMPLATE template1 ENCODING 'UTF8';

\c students

CREATE TABLE IF NOT EXISTS students (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    s_class TEXT,
    school TEXT,
    order_day SMALLINT,
    order_time TIME,
    order_cost SMALLINT
);
