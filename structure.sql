CREATE EXTENSION pgcrypto;

CREATE TABLE nodes (
    id TEXT PRIMARY KEY NOT NULL,
    id_short TEXT UNIQUE NOT NULL,
    name TEXT UNIQUE NOT NULL,
    private_key TEXT NOT NULL,
    public_key TEXT NOT NULL,
    host_key TEXT NOT NULL,
    ip TEXT NOT NULL,
    username TEXT NOT NULL,
    port TEXT NOT NULL,
);

CREATE TABLE tasks (
    id TEXT PRIMARY KEY NOT NULL,
    id_short TEXT UNIQUE NOT NULL,
    plugin TEXT NOT NULL,
    state INT NOT NULL,
    node TEXT NOT NULL,
    time TEXT,
    exit_code INT NOT NULL,
    log JSON,
    details JSONB NOT NULL
);

CREATE TABLE jobs (
    id TEXT PRIMARY KEY NOT NULL,
    id_short TEXT UNIQUE NOT NULL,
    name TEXT UNIQUE NOT NULL,
    cron TEXT,
    nodes JSONB NOT NULL,
    tasks JSON NOT NULL,
    task_history JSONB
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY NOT NULL,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);
