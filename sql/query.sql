-- name: GetRun :one
SELECT * FROM RunsData WHERE id = $1 LIMIT 1;

-- name: AddRunRaw :one
INSERT INTO RunsData (
    ascension_level,
    campfire_rested,
    campfire_upgraded,
    choose_seed,
    circlet_count,
    current_hp_per_floor,
    floor_reached,
    gold,
    gold_per_floor,
    is_beta, -- 10
    is_daily,
    is_endless,
    is_prod,
    is_trial,
    items_purchased_floors,
    items_purged_floors,
    local_time,
    max_hp_per_floor,
    neow_bonus,
    neow_cost, -- 20
    path_per_floor,
    path_taken,
    play_id,
    player_experience,
    playtime,
    potions_floor_spawned,
    potions_floor_usage,
    purchased_purges,
    score,
    seed_played, -- 30
    seed_source_timestamp,
    "timestamp",
    victory,
    win_rate,
    -- Default null/invalid values
    build_version, character_chosen
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
    $21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
    $31, $32, $33, $34,
    -- Default values
    1, 1
)
RETURNING RunsData.id;

-- name: SetRunText :exec
UPDATE RunsText SET
    build_version = $2,
    character_chosen = $3,
    items_purchased_names = $4,
    items_purged_names = $5
WHERE
    id = $1;

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

