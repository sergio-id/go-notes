START TRANSACTION;

CREATE SCHEMA IF NOT EXISTS "user";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "user"."user"
(
    id         uuid                     not null default (uuid_generate_v4()),
    email      varchar(64) UNIQUE       not null check ( email <> '' ),
    password   varchar(255)             not null check ( octet_length(password) <> 0 ),
    first_name varchar(64)              not null default '',
    last_name  varchar(64)              not null default '',
    updated_at timestamp with time zone not null default CURRENT_TIMESTAMP,
    created_at timestamp with time zone not null default NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS ix_user_id ON "user"."user" (id);

COMMIT;
