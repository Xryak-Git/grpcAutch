CREATE TABLE IF NOT EXISTS users
(
    id  INTEGER PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    pass_hash BLOB NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_email on users (email);

CREATE TABLE IF NOT EXISTS apps
(
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE
);