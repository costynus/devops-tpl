-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TYPE metric_types AS ENUM (
  'gauge',
  'counter'
);
CREATE TABLE IF NOT EXISTS public.metrics(
    id serial PRIMARY KEY,
    name VARCHAR(256),
    mtype metric_types,
    delta INTEGER,
    value DOUBLE PRECISION,
    hash VARCHAR(64),
    CONSTRAINT unique_name_indx UNIQUE (name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE public.metrics;
DROP TYPE metric_types;
-- +goose StatementEnd
