CREATE TABLE IF NOT EXISTS statuses (
    id SMALLSERIAL PRIMARY KEY,
    name VARCHAR(16) UNIQUE NOT NULL
);

CREATE INDEX idx_status_name_fulltext ON statuses USING gin (to_tsvector('russian', name));

CREATE TABLE IF NOT EXISTS tasks (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(64) NOT NULL,
    description TEXT NOT NULL,
    status_id SMALLINT NOT NULL,
    date TIMESTAMPTZ NOT NULL,
    deleted BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (status_id) REFERENCES statuses (id)
);

CREATE INDEX idx_tasks_status_id ON tasks (status_id);
CREATE INDEX idx_tasks_date ON tasks (date);

INSERT INTO statuses (name)
VALUES
    ('выполнено'),
    ('не выполнено');

