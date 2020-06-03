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
    data JSONB
);

create table jobs (
    id TEXT PRIMARY KEY NOT NULL,
    id_short TEXT UNIQUE NOT NULL,
    data JSONB
);
