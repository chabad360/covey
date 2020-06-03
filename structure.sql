CREATE TABLE nodes (
    id TEXT PRIMARY KEY NOT NULL,
    id_short TEXT UNIQUE NOT NULL,
    name TEXT UNIQUE NOT NULL,
    plugin TEXT NOT NULL,
    details JSONB NOT NULL
);

CREATE TABLE tasks (
    id TEXT PRIMARY KEY NOT NULL,
    id_short TEXT UNIQUE NOT NULL,
    plugin TEXT NOT NULL,
    state INT NOT NULL,
    node TEXT NOT NULL,
    time TEXT,
    log JSONB,
    details JSONB NOT NULL
);

create table jobs (
    id TEXT PRIMARY KEY NOT NULL,
    id_short TEXT UNIQUE NOT NULL,
    name TEXT UNIQUE NOT NULL,
    cron TEXT,
    nodes JSONB NOT NULL,
    tasks JSONB NOT NULL,
    task_history JSONB
);
