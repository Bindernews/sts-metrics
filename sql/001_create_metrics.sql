
-- Deduplicate names in other tables
CREATE TABLE StrCache(
    id SERIAL PRIMARY KEY,
    str TEXT NOT NULL,
    UNIQUE (str)
);
CREATE INDEX ON StrCache USING HASH(str);

CREATE OR REPLACE FUNCTION add_str(s TEXT) RETURNS INT
LANGUAGE SQL AS $$
    INSERT INTO StrCache(str) VALUES (s) ON CONFLICT (str) DO NOTHING;
    SELECT id FROM StrCache WHERE str = s;
$$;

CREATE OR REPLACE FUNCTION add_str_many(s TEXT[]) RETURNS TABLE(id INT)
LANGUAGE SQL AS $$
    INSERT INTO StrCache(str) SELECT unnest(s) ON CONFLICT (str) DO NOTHING;
    SELECT id FROM StrCache WHERE str = ANY(s);
$$;

CREATE OR REPLACE FUNCTION get_str(i INT) RETURNS TEXT
LANGUAGE SQL AS $$
    SELECT str FROM StrCache WHERE id = i;
$$;

CREATE OR REPLACE FUNCTION get_str_many(ids INT[]) RETURNS TEXT[]
LANGUAGE SQL AS $$
    SELECT ARRAY(SELECT str FROM StrCache WHERE id = ANY(ids));
$$;

CREATE TABLE RunsData
(
    id SERIAL PRIMARY KEY,
    ascension_level INT NOT NULL,
    /* REVERSE boss relics */
    build_version INT NOT NULL DEFAULT 1 REFERENCES StrCache(id),
    /* REVERSE campfire choices */
    campfire_rested INT,
    campfire_upgraded INT,
    /* REVERSE card choices */
    character_chosen INT NOT NULL DEFAULT 1 REFERENCES StrCache(id),
    choose_seed BOOLEAN NOT NULL,
    circlet_count INT,
    current_hp_per_floor INT ARRAY,
    /* REVERSE damage taken */
    -- EventChoices []RunSchemaJsonEventChoicesElem `json:"event_choices" yaml:"event_choices"`
    floor_reached INT NOT NULL ,
    gold INT NOT NULL ,
    gold_per_floor INT ARRAY,
    /* is ascention mode? if not, ascention_level will be 0 */
    is_beta BOOLEAN NOT NULL,
    is_daily BOOLEAN NOT NULL,
    is_endless BOOLEAN NOT NULL,
    is_prod BOOLEAN NOT NULL,
    is_trial BOOLEAN NOT NULL,
    items_purchased_floors INT ARRAY,
    items_purchased_ids INT ARRAY,
    -- Card remove floors
    items_purged_floors INT ARRAY,
    -- Card remove names
    items_purged_ids INT ARRAY,
    -- Encounter ID where player died
    killed_by INT NOT NULL DEFAULT 1 REFERENCES StrCache(id),
    -- Local time in YYYYmmddHHMMSS format
    local_time TEXT NOT NULL,
    master_deck INT ARRAY,
    -- Doctext
    max_hp_per_floor INT ARRAY,
    -- ID of player's Neow choice
    neow_bonus TEXT NOT NULL,
    -- TODO
    neow_cost TEXT NOT NULL,
    -- Path per floor as a single string
    path_per_floor TEXT NOT NULL,
    -- Path taken, stored as single string of characters
    path_taken TEXT NOT NULL,
    -- UUID for this run (UUID)
    play_id TEXT NOT NULL,
    -- XP gained at the end of the run
    player_experience INT NOT NULL ,
    -- Play time in seconds
    playtime INT NOT NULL,
    -- TODO Doc
    potions_floor_spawned INT ARRAY,
    -- Which floors the player used a potion on
    potions_floor_usage INT ARRAY,
    /* REVERSE PotionObtains */
    -- PurchasedPurges corresponds to the JSON schema field "purchased_purges".
    purchased_purges INT NOT NULL,
    /* REVERSE RelicObtains */
    -- Player's score at the end of the run
    score INT NOT NULL ,
    -- The run seed
    seed_played TEXT NOT NULL,
    -- TODO doc
    seed_source_timestamp INT,
    -- Timestamp corresponds to the JSON schema field "timestamp".
    "timestamp" TIMESTAMP,
    -- Victory corresponds to the JSON schema field "victory".
    victory BOOLEAN NOT NULL,
    -- WinRate corresponds to the JSON schema field "win_rate".
    win_rate FLOAT NOT NULL,
    UNIQUE (play_id)
);

CREATE OR REPLACE VIEW RunsText AS (
    SELECT
        R.id,
        (SELECT get_str(R.build_version))::text as build_version,
        (SELECT get_str(R.character_chosen))::text as character_chosen,
        (SELECT get_str_many(R.items_purchased_ids))::text[] as items_purchased_names,
        (SELECT get_str_many(R.items_purged_ids))::text[] as items_purged_names,
        (SELECT get_str(R.killed_by))::text as killed_by
    FROM RunsData AS R
);

CREATE OR REPLACE FUNCTION runstext_update() RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE RunsData SET
        build_version = (SELECT add_str(NEW.build_version)),
        character_chosen = (SELECT add_str(NEW.character_chosen)),
        items_purchased_ids = ARRAY((SELECT add_str_many(NEW.items_purchased_names))),
        items_purged_ids = ARRAY((SELECT add_str_many(NEW.items_purged_names))),
        killed_by = (SELECT add_str(NEW.killed_by))
    WHERE
        RunsData.id = NEW.id;
    RETURN NEW;
END $$;

CREATE TRIGGER runstext_update_trg INSTEAD OF UPDATE ON RunsText
    FOR EACH ROW EXECUTE FUNCTION runstext_update();


CREATE TABLE CampfireChoice(
    id serial primary key,
    run_id int not null references RunsData(id),
    cdata text,
    floor int not null,
    "key" int not null references StrCache(id)
);

CREATE TABLE DamageTaken(
    id serial primary key,
    run_id int not null references RunsData(id),
    enemies int not null references StrCache(id),
    floor int not null,
    turns int not null
);

CREATE TABLE BossRelics(
    id serial primary key,
    run_id int not null references RunsData(id),
    not_picked int array,
    picked int not null references StrCache(id)
);

CREATE TABLE CardChoices(
    id serial primary key,
    run_id int not null references RunsData(id),
    not_picked int array,
    picked int not null references StrCache(id),
    floor int not null
);

CREATE TABLE RelicObtains(
    id serial primary key,
    run_id int not null references RunsData(id),
    floor int not null,
    "key" int not null references StrCache(id)
);

CREATE TABLE PotionObtains(
    id serial primary key,
    run_id int not null references RunsData(id),
    floor int not null,
    "key" int not null references StrCache(id)
);

CREATE TABLE EventChoices(
    id serial primary key,
    run_id int not null references RunsData(id),
    -- damage_healed + damage_taken
	damage_delta int not null,
	event_name_id int not null references StrCache(id),
	floor int not null,
	-- Combo of gold gained+lost
	gold_delta int not null,
	-- Combo of max_hp gained+lost
	max_hp_delta int not null,
	player_choice_id int not null references StrCache(id),
	relics_obtained_ids int[]
);

-- Add empty string to cache
INSERT INTO StrCache(id, str) VALUES (1, '') ON CONFLICT DO NOTHING;

---- create above / drop below ----

drop trigger if exists runstext_update_trg ON RunsText;
drop function if exists runstext_update;
drop view if exists RunsText;
drop table if exists EventChoices;
drop table if exists BossRelics;
drop table if exists CardChoices;
drop table if exists DamageTaken;
drop table if exists CampfireChoice;
drop table if exists RelicObtains;
drop table if exists PotionObtains;
drop table if exists RunsData;
drop table if exists StrCache;
drop function if exists add_str;
drop function if exists add_str_many;
drop function if exists get_str;