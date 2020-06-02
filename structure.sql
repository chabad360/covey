create table nodes (
    id text primary key,
    id_short text unique,
    data jsonb
);

create table tasks (
    id text primary key,
    id_short text unique,
    data jsonb
);

create table jobs (
    id text primary key,
    id_short text unique,
    data jsonb
);
