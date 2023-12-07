CREATE TABLE IF NOT EXISTS statuses (
    id SMALLSERIAL PRIMARY KEY,
    name VARCHAR(16) UNIQUE NOT NULL
);

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

INSERT INTO statuses (name)
VALUES
    ('Выполнено'),
    ('Не выполнено');

