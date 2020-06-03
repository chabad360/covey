CREATE TABLE nodes (
    id TEXT PRIMARY KEY NOT NULL,
    id_short TEXT UNIQUE NOT NULL,
    name TEXT UNIQUE NOT NULL,
    plugin TEXT NOT NULL,
    details JSONB
);

CREATE TABLE tasks (
    id TEXT PRIMARY KEY NOT NULL,
    id_short TEXT UNIQUE NOT NULL,
    plugin TEXT NOT NULL,
    state INT NOT NULL,
    node TEXT NOT NULL,
    time TEXT,
    log JSONB,
    details JSONB
);

create table jobs (
    id TEXT PRIMARY KEY NOT NULL,
    id_short TEXT UNIQUE NOT NULL,
    data JSONB
);
