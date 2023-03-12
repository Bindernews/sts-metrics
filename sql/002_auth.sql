
CREATE SCHEMA IF NOT EXISTS auth;
SET search_path TO auth, public;

CREATE TABLE IF NOT EXISTS users(
    id serial primary key ,
    email text not null unique
);

CREATE TABLE IF NOT EXISTS scopes(
    id serial primary key,
    "key" text not null, -- code ID for scope
    "desc" text not null, -- long description
    parent int references scopes(id) -- parent scope, if applicable
);

CREATE TABLE IF NOT EXISTS users_to_scopes(
    user_id int not null references users(id) on delete cascade ,
    scope_id int not null references scopes(id) on delete cascade ,
    primary key (user_id, scope_id)
);

CREATE OR REPLACE FUNCTION user_add(email_ text) RETURNS int
LANGUAGE SQL AS $$
    INSERT INTO users (email)
    VALUES (email_)
    RETURNING users.id;
$$;

CREATE OR REPLACE FUNCTION user_set_scopes(email_ text, scope_list text[]) RETURNS void
LANGUAGE plpgsql AS $$
DECLARE
    uid int;
    scope_ids int[];
BEGIN
    PERFORM set_config('search_path', 'auth', true);
    SELECT id INTO uid FROM users WHERE email = email_;
    if not FOUND then
        RAISE EXCEPTION 'user % not found', email_;
    end if;
    SELECT array(SELECT id FROM scopes WHERE key = ANY(scope_list)) INTO scope_ids;
    if array_length(scope_ids, 1) != array_length(scope_list, 1) then
        RAISE EXCEPTION 'one or more scopes of % are invalid', scope_list;
    end if;
    DELETE FROM users_to_scopes WHERE user_id = uid;
    INSERT INTO users_to_scopes(user_id, scope_id) SELECT uid, unnest(scope_ids);
END $$;

CREATE OR REPLACE FUNCTION user_has_scopes(email_ text, scope_list text[]) RETURNS bool
LANGUAGE SQL AS $$
    SELECT bool_and(scope_id is not null)
    FROM auth.scopes s
        CROSS JOIN (SELECT id FROM auth.users WHERE email = email_) uid
        RIGHT JOIN (SELECT unnest(scope_list) t) txt ON txt.t = s.key
        LEFT JOIN auth.users_to_scopes uts ON s.id = uts.scope_id AND uid.id = uts.user_id
$$;

-- TEST DATA
INSERT INTO scopes(key, "desc")
VALUES
    ('stats:view', 'Ability to view statistics'),
    ('getrun', 'Retrieve runs as JSON data'),
    ('admin', 'Access the admin interface');

---- create above / drop below ----

set search_path to auth,public;
drop table if exists users_to_scopes;
drop table if exists scopes;
drop table if exists users;
drop function if exists user_add;
drop function if exists user_set_scopes;
drop function if exists user_has_scopes;
drop schema if exists auth;
