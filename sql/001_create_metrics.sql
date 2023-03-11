
-- Deduplicate names in other tables
CREATE TABLE StrCache(
    id int primary key generated by default as identity,
    str text not null,
    unique (str)
);

-- Add empty string to cache
INSERT INTO StrCache(str) VALUES ('');

CREATE OR REPLACE FUNCTION str_cache_add(s text[]) RETURNS void
LANGUAGE SQL AS $$
    WITH newstr AS (SELECT t.t FROM unnest(s) t JOIN StrCache ON t != StrCache.str)
    INSERT INTO StrCache(str) SELECT t FROM newstr
    ON CONFLICT DO NOTHING;
$$;

CREATE OR REPLACE FUNCTION str_cache_to_id(s text[]) RETURNS TABLE(id int)
LANGUAGE SQL AS $$
    WITH str_list AS (SELECT row_number() over () as ix, t FROM (SELECT unnest(s) as t) list)
    SELECT S.id FROM StrCache S RIGHT JOIN str_list ON S.str = str_list.t ORDER BY str_list.ix;
$$;

CREATE OR REPLACE FUNCTION str_cache_to_str(ids int[]) RETURNS TABLE(t text)
LANGUAGE SQL AS $$
    WITH id_list AS (SELECT row_number() over () as ix, id FROM (SELECT unnest(ids) as id) list)
    SELECT S.str as t FROM StrCache S RIGHT JOIN id_list ON S.id = id_list.id ORDER BY id_list.ix;
$$;

CREATE OR REPLACE FUNCTION get_str(i INT) RETURNS TEXT
LANGUAGE SQL AS $$
    SELECT str FROM StrCache WHERE id = i;
$$;

-- A card and its upgrade count. This is so that we can organize things by base card name
-- while also being able to differentiate by upgrade if necessary.
CREATE TABLE CardSpecs(
    id int primary key generated by default as identity,
    -- Card name WITHOUT upgrades
    card text not null,
    -- How many upgrades this card has
    upgrades int not null,
    unique (card, upgrades)
);

CREATE INDEX ON CardSpecs USING btree (card);

-- Card spec joined with strings.
CREATE VIEW CardSpecsEx AS
(
    WITH c AS (SELECT c.*,
                      (CASE WHEN c.upgrades > 0 THEN '+' || c.upgrades END)::text as suffix
               FROM CardSpecs c)
    SELECT c.*, concat(c.card, c.suffix) as card_full
    FROM c
);


CREATE TYPE card_spec_io AS (card text, upg int);

CREATE FUNCTION card_spec_add(cards_in card_spec_io[]) RETURNS void
LANGUAGE SQL AS $$
    WITH cards AS (SELECT c.card, c.upg FROM unnest(cards_in) c)
    INSERT INTO CardSpecs (card, upgrades)
    SELECT * FROM cards cs
    WHERE NOT exists(SELECT id FROM CardSpecs c WHERE c.card = cs.card AND c.upgrades = cs.upg)
$$;

CREATE FUNCTION card_spec_to_id(cards_in card_spec_io[]) RETURNS TABLE(id int)
LANGUAGE SQL AS $$
    WITH src AS (SELECT row_number() over () as ix, c.card, c.upg FROM (SELECT * FROM unnest(cards_in)) c)
    SELECT cs.id
    FROM cardspecs cs
    JOIN src ON src.card = cs.card AND src.upg = cs.upgrades
$$;

CREATE TABLE RunsData(
    id int primary key generated by default as identity,
    ascension_level int not null,
    /* REVERSE boss relics */
    build_version int not null default 1 references StrCache(id),
    /* REVERSE campfire choices */
    campfire_rested int not null default 0,
    campfire_upgraded int not null default 0,
    /* REVERSE card choices */
    character_id int not null default 1 references StrCache(id),
    choose_seed boolean not null,
    circlet_count int not null default 0,
    /* REVERSE damage taken */
    -- EventChoices []RunSchemaJsonEventChoicesElem `json:"event_choices" yaml:"event_choices"`
    floor_reached int not null ,
    gold int not null ,
    -- Encounter ID where player died
    killed_by int not null default 1 references StrCache(id),
    -- Local time in YYYYmmddHHMMSS format
    local_time text not null,
    /** REVERSE master deck */
    -- ID of player's Neow choice
    neow_bonus_id int not null default 1 references StrCache(id),
    -- Trade-off when making Neow choice
    neow_cost_id int not null default 1 references StrCache(id),
    -- Path per floor as a single string
    path_per_floor text not null,
    -- Path taken, stored as single string of characters
    path_taken text not null,
    -- UUID for this run (UUID)
    play_id text not null,
    -- XP gained at the end of the run
    player_experience int not null ,
    -- Play time in seconds
    playtime int not null,
    /* REVERSE PotionObtains */
    -- How many purges were purchased from the shop
    purchased_purges int not null,
    /* REVERSE RelicObtains */
    -- Player's score at the end of the run
    score int not null ,
    -- The run seed
    seed_played text not null,
    -- TODO doc
    seed_source_timestamp INT,
    -- Special seed if given, 0 if not
    special_seed int not null default 0,
    -- Timestamp corresponds to the JSON schema field "timestamp".
    "timestamp" TIMESTAMP,
    -- Victory corresponds to the JSON schema field "victory".
    victory boolean not null,
    -- WinRate corresponds to the JSON schema field "win_rate".
    win_rate float not null,
    UNIQUE (play_id)
);
CREATE INDEX runsdata_character_index ON RunsData (character_id);

CREATE TYPE flag_kind AS ENUM (
    'ascension', 'beta', 'daily', 'endless', 'prod', 'trial'
);

CREATE TABLE RunFlags(
    run_id int not null references RunsData(id),
    flag flag_kind,
    PRIMARY KEY (run_id, flag)
);

CREATE TABLE PerFloorData(
    run_id int not null references RunsData(id),
    floor int2 not null,
    gold int not null,
    current_hp int not null,
    max_hp int not null,
    primary key (run_id, floor)
);

-- Split some array data from the main table.
-- This reduces tuple size, and de-clutters some of the insert code.
CREATE TABLE RunArrays(
    run_id int not null references RunsData(id),
    -- List of daily mods
    daily_mods int array,
    -- Card purchase floors
    items_purchased_floors int array,
    -- Card purchase IDs
    items_purchased_ids int array,
    -- Card remove floors
    items_purged_floors int array,
    -- Card remove IDs
    items_purged_ids int array,
    -- Which floors potions were available as options
    potions_floor_spawned int array,
    -- Which floors the player used a potion on
    potions_floor_usage int array,
    -- List of relic names
    relic_ids int array,
    primary key (run_id)
);

-- Stores any extra fields not parsed normally
CREATE TABLE runs_extra(
    run_id int not null references RunsData(id),
    extra jsonb not null,
    primary key (run_id)
);

CREATE TABLE CampfireChoice(
    id int primary key generated by default as identity,
    run_id int not null references RunsData(id),
    "data" int references StrCache(id),
    floor int not null,
    "key" int not null references StrCache(id)
);

CREATE TABLE DamageTaken(
    id int primary key generated by default as identity,
    run_id int not null references RunsData(id),
    enemies int not null references StrCache(id),
    damage float4 not null,
    floor int not null,
    turns int not null
);

CREATE TABLE BossRelics(
    id int primary key generated by default as identity,
    run_id int not null references RunsData(id),
    not_picked int array,
    picked int not null references StrCache(id),
    -- order
    ord int2 not null
);

CREATE TABLE CardChoices(
    id int primary key generated by default as identity,
    run_id int not null references RunsData(id),
    not_picked int array,
    picked int not null references CardSpecs(id),
    floor int not null
);

-- CREATE TABLE CardChoicesNotPicked(
--     id int primary key generated by default as identity,
--     choice_id int not null references CardChoices(id),
--     card_id int not null references CardSpecs(id)
-- );

CREATE TABLE RelicObtains(
    id int primary key generated by default as identity,
    run_id int not null references RunsData(id),
    floor int2 not null,
    "key" int not null references StrCache(id)
);

CREATE TABLE PotionObtains(
    id int primary key generated by default as identity,
    run_id int not null references RunsData(id),
    floor int2 not null,
    "key" int not null references StrCache(id)
);

CREATE TABLE EventChoices(
    id int primary key generated by default as identity,
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

CREATE TABLE MasterDecks(
    id int primary key generated by default as identity,
    run_id int not null references RunsData(id),
    card_id int not null references CardSpecs(id),
    -- Copies of this card in the deck
    count int2 not null,
    -- These should be unique
    UNIQUE (run_id, card_id, count)
);

-- Storage for compressed original JSON, in case we need to re-parse it.
-- Periodically this table will be cleared and the contents will be
-- compressed together into one larger archive file.
--
-- The status is there in case something goes wrong during an export,
-- or if two different exports run at the same time for some reason.
CREATE TABLE RawJsonArchive(
    id int primary key generated by default as identity,
    bdata json not null,
    -- The play_id of the data
    play_id text not null,
    -- 0 = not exported, 1 = export in progress, 2 = export done
    status int2 not null default 0,
    unique (play_id)
);

-- Materialized so we don't have scan RunsData for the list of unique characters.
CREATE MATERIALIZED VIEW IF NOT EXISTS character_list AS (
    SELECT DISTINCT S.id, S.str as "name" FROM StrCache S
    LEFT JOIN RunsData R ON S.id = R.character_id
    WHERE S.id = R.character_id
    ORDER BY S.str
);

-- Auto-update whenever we add a new run.
CREATE FUNCTION character_list_refresh() RETURNS TRIGGER
LANGUAGE plpgsql AS $$
BEGIN
    REFRESH MATERIALIZED VIEW character_list;
    RETURN NULL;
END $$;

CREATE TRIGGER character_list_refresh AFTER INSERT OR UPDATE OR DELETE
ON RunsData FOR EACH STATEMENT EXECUTE FUNCTION character_list_refresh();


-- run_to_json: Does most of the work to turn tables back into JSON.
-- Go (or other clients) transform path_taken and path_per_floor, then merge the transformed
-- values, raw, and extra, into one big json blob. The result should be (if keys are sorted)
-- identical to the input json.
CREATE OR REPLACE FUNCTION run_to_json(runid int)
    RETURNS TABLE(raw json, path_per_floor text, path_taken text, extra json)
LANGUAGE SQL AS $$
    WITH
        rn AS (
            SELECT r.*,
                   get_str(r.build_version) as s_build_version,
                   get_str(r.character_id) as s_character,
                   get_str(r.killed_by) as s_killed_by,
                   get_str(r.neow_bonus_id) as s_neow_bonus,
                   get_str(r.neow_cost_id) as s_neow_cost,
                   re.extra as extra
            FROM runsdata r LEFT JOIN runs_extra re on r.id = re.run_id
            WHERE r.id = runid
        ),
        dmg AS (
            SELECT d.floor, d.turns, s.str as enemies
            FROM damagetaken d JOIN strcache s on s.id = d.enemies
            WHERE d.run_id = runid
        ),
        -- Generate deck strings
        deck AS (
            WITH m AS (SELECT * FROM masterdecks m WHERE m.run_id = runid),
                 ex AS (SELECT generate_series(1, (SELECT max(m.count))) v FROM m)
            SELECT concat(s.card, s.suffix) AS card
            FROM m LEFT JOIN ex ON ex.v <= m.count JOIN CardSpecsEx s ON s.id = m.card_id
        ),
        card_choices AS (
            SELECT s.card_full as picked,
                   c.floor,
                   array(SELECT s.card_full
                         FROM CardSpecsEx s
                         WHERE s.id = any(c.not_picked)
                       ) as not_picked
            FROM cardchoices c JOIN cardspecsex s on c.picked = s.id
            WHERE c.run_id = runid
        ),
        per_floor AS (
            SELECT array_agg(pf.gold ORDER BY pf.floor) as gold,
                   array_agg(pf.current_hp ORDER BY pf.floor) as current_hp,
                   array_agg(pf.max_hp ORDER BY pf.floor) as max_hp
            FROM perfloordata pf
            WHERE pf.run_id = runid
        ),
        potions_obt AS (
            SELECT p.floor, s.str as "key"
            FROM potionobtains p JOIN strcache s on p.key = s.id
            WHERE p.run_id = runid
        ),
        relics_obt AS (
            SELECT r.floor, s.str as "key"
            FROM relicobtains r JOIN strcache s on r.key = s.id
            WHERE r.run_id = runid
        ),
        campfires AS (
            SELECT c.floor,
                   get_str(c.key),
                   (SELECT s.str FROM strcache s WHERE c.data = s.id) as data
            FROM (SELECT * FROM campfirechoice c WHERE c.run_id = runid) c
        ),
        events AS (
            SELECT e.floor,
                   get_str(e.event_name_id) as event_name,
                   get_str(e.player_choice_id) as player_choice,
                   (CASE WHEN e.damage_delta > 0 THEN e.damage_delta END) as damage_healed,
                   (CASE WHEN e.damage_delta < 0 THEN -e.damage_delta END) as damage_taken,
                   (CASE WHEN e.gold_delta > 0 THEN e.gold_delta END) as gold_gain,
                   (CASE WHEN e.gold_delta < 0 THEN -e.gold_delta END) as gold_loss,
                   (CASE WHEN e.max_hp_delta > 0 THEN e.max_hp_delta END) as max_hp_gain,
                   (CASE WHEN e.max_hp_delta < 0 THEN -e.max_hp_delta END) as max_hp_loss,
                   array(SELECT str_cache_to_str(e.relics_obtained_ids)) as relics_obtained
            FROM eventchoices e
            WHERE e.run_id = runid
        ),
        boss_relics AS (
            SELECT get_str(br.picked) as picked,
                   array(SELECT str_cache_to_str(br.not_picked)) as not_picked
            FROM bossrelics br
            WHERE br.run_id = runid
            ORDER BY br.ord
        ),
        arrs AS (
            SELECT array(SELECT str_cache_to_str(a.daily_mods)) as daily_mods,
                   a.items_purchased_floors,
                   array(SELECT str_cache_to_str(a.items_purchased_ids)) as items_purchased,
                   a.items_purged_floors,
                   array(SELECT str_cache_to_str(a.items_purged_ids)) as items_purged,
                   a.potions_floor_spawned,
                   a.potions_floor_usage,
                   array(SELECT str_cache_to_str(a.relic_ids)) as relics
            FROM runarrays a
            WHERE a.run_id = runid
        ),
        flags AS (SELECT * FROM runflags WHERE run_id = runid)
    SELECT json_build_object(
        'ascension_level', rn.ascension_level,
        'boss_relics', (SELECT json_agg(row_to_json(boss_relics)) FROM boss_relics),
        'build_version', rn.s_build_version,
        'campfire_choices', (SELECT json_agg(json_strip_nulls(row_to_json(campfires))) FROM campfires),
        'campfire_rested', rn.campfire_rested,
        'campfire_upgraded', rn.campfire_upgraded,
        'card_choices', (SELECT json_agg(row_to_json(card_choices)) FROM card_choices),
        'character_chosen', rn.s_character,
        'chose_seed', rn.choose_seed,
        'circlet_count', rn.circlet_count,
        'current_hp_per_floor', (SELECT array_to_json(per_floor.current_hp) FROM per_floor),
        'daily_mods', (SELECT array_to_json(arrs.daily_mods) FROM arrs),
        'damage_taken', (SELECT json_agg(row_to_json(dmg)) FROM dmg),
        'event_choices', (SELECT json_agg(json_strip_nulls(row_to_json(events))) FROM events),
        'floor_reached', rn.floor_reached,
        'gold', rn.gold,
        'gold_per_floor', (SELECT array_to_json(per_floor.gold) FROM per_floor),
        'is_ascension', (SELECT count(*) FROM flags WHERE flags.flag = 'ascension'),
        'is_beta', (SELECT count(*)>0 FROM flags WHERE flags.flag = 'beta'),
        'is_daily', (SELECT count(*)>0 FROM flags WHERE flags.flag = 'daily'),
        'is_endless', (SELECT count(*)>0 FROM flags WHERE flags.flag = 'endless'),
        'is_prod', (SELECT count(*)>0 FROM flags WHERE flags.flag = 'prod'),
        'is_trial', (SELECT count(*)>0 FROM flags WHERE flags.flag = 'trial'),
        'item_purchase_floors', (SELECT array_to_json(arrs.items_purged_floors) FROM arrs),
        'items_purchased', (SELECT array_to_json(arrs.items_purchased) FROM arrs),
        'items_purged_floors', (SELECT array_to_json(arrs.items_purged_floors) FROM arrs),
        'items_purged', (SELECT array_to_json(arrs.items_purged) FROM arrs),
        'killed_by', rn.s_killed_by,
        'local_time', rn.local_time,
        'master_deck', (SELECT json_agg(deck.card) FROM deck),
        'max_hp_per_floor', (SELECT array_to_json(per_floor.max_hp) FROM per_floor),
        'neow_bonus', rn.s_neow_bonus,
        'neow_cost', rn.s_neow_cost,
        'play_id', rn.play_id,
        'player_experience', rn.player_experience,
        'playtime', rn.playtime,
        'potions_floor_spawned', (SELECT array_to_json(arrs.potions_floor_spawned) FROM arrs),
        'potions_floor_usage', (SELECT array_to_json(arrs.potions_floor_usage) FROM arrs),
        'potions_obtained', (SELECT json_agg(row_to_json(potions_obt)) FROM potions_obt),
        'purchased_purges', rn.purchased_purges,
        'relics', (SELECT array_to_json(arrs.relics) FROM arrs),
        'relics_obtained', (SELECT json_agg(row_to_json(relics_obt)) FROM relics_obt),
        'score', rn.score,
        'seed_played', rn.seed_played,
        'seed_source_timestamp', rn.seed_source_timestamp,
        'special_seed', rn.special_seed,
        'timestamp', rn."timestamp",
        'victory', rn.victory,
        'win_rate', rn.win_rate
    ) as raw,
        rn.path_per_floor,
        rn.path_taken,
        rn.extra::json
    FROM rn;
$$;

---- create above / drop below ----

drop trigger if exists character_list_refresh ON RunsData;
drop function if exists character_list_refresh;
drop materialized view if exists character_list;

drop table if exists MasterDecks;
drop table if exists EventChoices;
drop table if exists BossRelics;
drop table if exists CardChoices;
drop table if exists DamageTaken;
drop table if exists CampfireChoice;
drop table if exists RelicObtains;
drop table if exists PotionObtains;
drop table if exists RunFlags;
drop table if exists PerFloorData;
drop table if exists RunArrays;
drop table if exists runs_extra;
drop table if exists rawjsonarchive;
drop view if exists CardSpecsEx;
drop table if exists CardSpecs;
drop table if exists RunsData;
drop table if exists StrCache;
drop type if exists flag_kind;
drop function if exists run_to_json;
drop function if exists get_str;
drop function if exists str_cache_add;
drop function if exists str_cache_to_id;
drop function if exists str_cache_to_str;
drop function if exists card_spec_add;
drop function if exists card_spec_to_id;
drop type if exists card_spec_io;
