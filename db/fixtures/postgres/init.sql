START TRANSACTION;

--------------------------------------------schema--------------------------------------------

---------------------------------------------user---------------------------------------------

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

---------------------------------------------category---------------------------------------------

CREATE SCHEMA IF NOT EXISTS category;

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

---------------------------------------------note---------------------------------------------

CREATE SCHEMA IF NOT EXISTS note;

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

--------------------------------------------data--------------------------------------------

INSERT INTO "user"."user" (id, email, password, first_name, last_name)
VALUES ('ca9091f2-23e0-42ab-9ef5-70a043757e73', 'admin@localhost', '$argon2id$v=19$m=65536,t=1,p=4$q1/RKAjx24rSnkYerSiDiw$a0m79mCR6N5jpRmQkIE6esj3cULCa+R+tymD95Fu4hU', 'admin@localhost', '12345678');

INSERT INTO category.category (id, user_id, title, pinned, priority)
VALUES ('6ed6253a-6101-46fd-accb-6726a2a44024', 'ca9091f2-23e0-42ab-9ef5-70a043757e73', 'category name', false, 0);

INSERT INTO note.note (user_id, category_id, title, content, pinned, priority)
VALUES ('ca9091f2-23e0-42ab-9ef5-70a043757e73', '6ed6253a-6101-46fd-accb-6726a2a44024', 'note name', 'content data',
        true, 10);

COMMIT;
