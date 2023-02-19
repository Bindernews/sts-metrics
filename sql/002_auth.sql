
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
    user_id int not null references users(id),
    scope_id int not null references scopes(id),
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
BEGIN
    SELECT id INTO uid FROM users WHERE email = email_;
    if not FOUND then
        RAISE EXCEPTION 'user % not found', email_;
    end if;
    CREATE TEMPORARY TABLE scope_ids ON COMMIT DROP AS
        SELECT id FROM scopes WHERE key = ANY(scope_list);
    if (SELECT count(id) FROM scope_ids) != array_length(scope_list, 1) then
        RAISE EXCEPTION 'one or more scopes of % are invalid', scope_list;
    end if;
    DELETE FROM users_to_scopes WHERE user_id = uid;
    INSERT INTO users_to_scopes(user_id, scope_id) SELECT uid, scope_ids.id FROM scope_ids;
END $$;

CREATE OR REPLACE FUNCTION user_has_scopes(emaila text, scope_list text[]) RETURNS int
LANGUAGE SQL AS $$
    WITH
        usc AS (
            SELECT scope_id FROM users_to_scopes
            LEFT JOIN users u on u.id = users_to_scopes.user_id
            WHERE u.email = emaila
        )
    SELECT count(usc.scope_id) FROM usc
    LEFT JOIN scopes ON usc.scope_id = scopes.id
    WHERE scopes.key = ANY(scope_list);
$$;

-- TEST DATA
-- SELECT user_add({{ .test_user }});
-- insert into scopes(key, "desc") values ('stats:view', 'Ability to view statistics');
-- SELECT user_set_scopes({{ .test_user }}, '{stats:view}');

---- create above / drop below ----

drop table if exists users_to_scopes;
drop table if exists scopes;
drop table if exists users;
drop function if exists user_add;
drop function if exists user_set_scopes;
drop function if exists user_has_scopes;
