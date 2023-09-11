START TRANSACTION;

CREATE SCHEMA IF NOT EXISTS note;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS note.note
(
    id          uuid                                                 not null default (uuid_generate_v4()),
    user_id     uuid references "user"."user" (id) on delete cascade not null,
    category_id uuid references category.category (id) on delete cascade      default null,
    title       varchar(64)                                          not null check (title <> ''),
    content     text                                                 not null default '',
    pinned      boolean                                              not null default false,
    priority    integer                                              not null default 0,
    updated_at  timestamp with time zone                             not null default CURRENT_TIMESTAMP,
    created_at  timestamp with time zone                             not null default now()
);

CREATE UNIQUE INDEX IF NOT EXISTS ix_notes_id ON note.note (id);
CREATE INDEX IF NOT EXISTS ix_notes_category_id ON note.note (category_id);

COMMIT;
