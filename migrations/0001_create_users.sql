CREATE TABLE users (
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username        TEXT NOT NULL UNIQUE CHECK (username = LOWER(username)),
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    email           TEXT NOT NULL,
    password_hash   TEXT NOT NULL
);