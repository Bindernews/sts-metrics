-- name: GetRun :one
SELECT * FROM Runs WHERE id = $1 LIMIT 1;

-- name: AddStr :one
SELECT add_str($1);

-- name: AddRun :one
WITH
    build_version AS (SELECT add_str($2)),
    character_chosen AS (SELECT add_str($5)),
    items_purchased AS (SELECT add_str_many($18::text[])),
    items_purged AS (SELECT add_str_many($20::text[])),
    killed_by_id AS (SELECT add_str($21))
INSERT INTO Runs (
    ascension_level,
    build_version,
    campfire_rested,
    campfire_upgraded,
    character_chosen,
    choose_seed,
    circlet_count,
    current_hp_per_floor,
    floor_reached,
    gold,
    gold_per_floor,
    is_beta,
    is_daily,
    is_endless,
    is_prod,
    is_trial,
    items_purchased_floors,
    items_purchased_ids,
    items_purged_floors,
    items_purged_ids,
    killed_by,
    local_time,
    max_hp_per_floor,
    neow_bonus,
    neow_cost, -- 25
    -- TODO path per floor,
    -- TODO path taken,
    play_id,
    player_experience,
    playtime,
    potions_floor_spawned,
    potions_floor_usage, -- 30
    purchased_purges,
    score,
    seed_played,
    seed_source_timestamp,
    ctimestamp, -- 35
    victory,
    win_rate
) VALUES (
    $1, build_version, $3, $4, character_chosen,
    $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15,
    $16, $17, items_purchased, $19, items_purged,
    killed_by_id, $22, $23, $24, $25,
    $26, $27, $28, $29, $30,
    $31, $32, $33, $34, $35,
    $36, $37
)
RETURNING Runs.id;

-- name: AddCampfire :one
INSERT INTO CampfireChoice (run_id, cdata, floor, "key")
    VALUES ($1, $2, $3, (SELECT add_str(@ckey)))
    RETURNING CampfireChoice.id;

-- name: AddDamageTaken :one
INSERT INTO DamageTaken (run_id, enemies, floor, turns)
    VALUES ($1, (SELECT add_str(@enemies)), $2, $3)
    RETURNING DamageTaken.id;

-- name: AddCardChoice :one
INSERT INTO CardChoices (run_id, floor, not_picked, picked)
    VALUES ($1, $2, (SELECT add_str_many(@not_picked::text[])), (SELECT add_str(@picked)))
    RETURNING CardChoices.id;

-- name: AddRelicObtain :exec
INSERT INTO RelicObtains (run_id, floor, "key")
    VALUES ($1, $2, (SELECT add_str(@ckey::text)));

-- name: AddPotionObtain :exec
INSERT INTO PotionObtains (run_id, floor, "key")
    VALUES ($1, $2, (SELECT add_str(@ckey::text)));

-- name: GetCampfires :many
SELECT CC.id, CC.cdata, CC.floor, StrCache.str as "key" FROM CampfireChoice AS CC
    LEFT JOIN StrCache ON CC.key = StrCache.id
    WHERE CC.id = $1
    ORDER BY floor;

