START TRANSACTION;

CREATE SCHEMA IF NOT EXISTS category;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS category.category
(
    id         uuid                                                 not null default (uuid_generate_v4()),
    user_id    uuid references "user"."user" (id) on delete cascade not null,
    title      varchar(64)                                          not null check (title <> ''),
    pinned     boolean                                              not null default false,
    priority   integer                                              not null default 0,
    updated_at timestamp with time zone                             not null default CURRENT_TIMESTAMP,
    created_at timestamp with time zone                             not null default now()
);

CREATE UNIQUE INDEX IF NOT EXISTS ix_category_id ON category.category (id);

COMMIT;
