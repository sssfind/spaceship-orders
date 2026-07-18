-- +goose Up
CREATE TABLE IF NOT EXISTS users (
                                     id BIGSERIAL PRIMARY KEY,
                                     uuid UUID DEFAULT uuid_generate_v4() UNIQUE NOT NULL,
    login VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    notification_methods JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
                             );

-- Индекс для ускорения выборки по логину при авторизации
CREATE INDEX IF NOT EXISTS idx_users_login ON users(login);

-- +goose Down
DROP TABLE IF EXISTS users;