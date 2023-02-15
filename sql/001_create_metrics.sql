
-- Deduplicate names in other tables
CREATE TABLE StrCache
(
    id SERIAL PRIMARY KEY,
    str TEXT NOT NULL
);
CREATE INDEX ON StrCache USING HASH(str);

CREATE TABLE Runs (
    id SERIAL PRIMARY KEY,
    ascension_level INT NOT NULL,
    /* REVERSE boss relics */
    build_version INT NOT NULL REFERENCES StrCache(id),
    /* REVERSE campfire choices */
    campfire_rested INT,
    campfire_upgraded INT,
    /* REVERSE card choices */
    character_chosen INT NOT NULL REFERENCES StrCache(id),
    choose_seed BOOLEAN NOT NULL,
    circlet_count INT,
    current_hp_per_floor INT ARRAY,
    /* REVERSE damage taken */
    -- EventChoices []RunSchemaJsonEventChoicesElem `json:"event_choices" yaml:"event_choices"`
    floor_reached INT,
    gold INT,
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
    killed_by INT REFERENCES StrCache(id),
    -- Local time in YYYYmmddHHMMSS format
    local_time TEXT NOT NULL,
    /* REVERSE master deck */
    -- Doctext
    max_hp_per_floor INT ARRAY,
    -- ID of player's Neow choice
    neow_bonus TEXT NOT NULL,
    -- TODO
    neow_cost TEXT NOT NULL,
    /*
    // PathPerFloor corresponds to the JSON schema field "path_per_floor".
    PathPerFloor []RunSchemaJsonPathPerFloorElem `json:"path_per_floor" yaml:"path_per_floor"`

    // PathTaken corresponds to the JSON schema field "path_taken".
    PathTaken []string `json:"path_taken" yaml:"path_taken"`
    */
    -- UUID for this run (UUID)
    play_id TEXT NOT NULL,
    -- XP gained at the end of the run
    player_experience INT,
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
    score INT,
    -- The run seed
    seed_played TEXT NOT NULL,
    -- TODO doc
    seed_source_timestamp INT,
    -- Timestamp corresponds to the JSON schema field "timestamp".
    ctimestamp TIMESTAMP,
    -- Victory corresponds to the JSON schema field "victory".
    victory BOOLEAN,
    -- WinRate corresponds to the JSON schema field "win_rate".
    win_rate INT
);

CREATE TABLE CampfireChoice
(
    id SERIAL PRIMARY KEY,
    run_id INT NOT NULL REFERENCES Runs(id),
    cdata TEXT,
    floor INT NOT NULL,
    "key" INT NOT NULL REFERENCES StrCache(id)
);

CREATE TABLE DamageTaken
(
    id SERIAL PRIMARY KEY,
    run_id INT NOT NULL REFERENCES Runs(id),
    enemies INT NOT NULL REFERENCES StrCache(id),
    floor INT NOT NULL,
    turns INT NOT NULL
);

CREATE TABLE BossRelics
(
    id SERIAL PRIMARY KEY,
    not_picked INT ARRAY,
    picked INT NOT NULL REFERENCES StrCache(id)
);

CREATE TABLE CardChoices
(
    id SERIAL PRIMARY KEY,
    run_id INT NOT NULL REFERENCES Runs(id),
    not_picked INT ARRAY,
    picked INT NOT NULL REFERENCES StrCache(id),
    floor INT NOT NULL
);

Create TABLE RunsToBossRelics
(
    run_id INT NOT NULL REFERENCES Runs(id),
    relic_id INT NOT NULL REFERENCES BossRelics(id),
    aindex INT
);

CREATE TABLE RelicObtains
(
    id SERIAL PRIMARY KEY,
    run_id INT NOT NULL REFERENCES Runs(id),
    floor INT NOT NULL,
    "key" INT REFERENCES StrCache(id)
);

CREATE TABLE PotionObtains
(
    id SERIAL PRIMARY KEY,
    run_id INT NOT NULL REFERENCES Runs(id),
    floor INT NOT NULL,
    "key" INT NOT NULL REFERENCES StrCache(id)
);


CREATE FUNCTION add_str(s TEXT) RETURNS INT
LANGUAGE SQL AS $$
    INSERT INTO StrCache(str) VALUES (s)
    ON CONFLICT (str) DO NOTHING
    RETURNING (SELECT id FROM StrCache WHERE str = s);
$$;

CREATE OR REPLACE FUNCTION add_str_many(s TEXT[]) RETURNS TABLE(id INT)
LANGUAGE SQL AS $$
    WITH
        elems(e) AS (SELECT unnest(s))
    INSERT INTO StrCache(str) SELECT * FROM elems
    ON CONFLICT (str) DO NOTHING
    RETURNING (SELECT id FROM StrCache LEFT JOIN elems ON str = e);
$$;


---- create above / drop below ----

drop table if exists BossRelics;
drop table if exists CardChoices;
drop table if exists DamageTaken;
drop table if exists CampfireChoice;
drop table if exists RelicObtains;
drop table if exists PotionObtains;
drop table if exists Runs;
drop table if exists StrCache;


