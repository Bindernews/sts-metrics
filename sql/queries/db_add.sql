-- name: GetRun :one
SELECT * FROM RunsData WHERE id = $1 LIMIT 1;

-- name: AddRunRaw :one
INSERT INTO RunsData
    (ascension_level, campfire_rested, campfire_upgraded,
     choose_seed, circlet_count, current_hp_per_floor, floor_reached, gold, gold_per_floor,
     items_purchased_floors, items_purged_floors, local_time, max_hp_per_floor, neow_bonus,
     neow_cost, path_per_floor, path_taken, play_id, player_experience, playtime,
     potions_floor_spawned, potions_floor_usage, purchased_purges, score, seed_played,
     seed_source_timestamp, "timestamp", victory, win_rate)
VALUES (
    $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,
    $21,$22,$23,$24,$25,$26,$27,$28,$29
)
RETURNING RunsData.id;

-- name: DoesRunExist :one
SELECT count(id)::boolean FROM RunsData R WHERE R.play_id = $1;

-- name: AddRunText :exec
UPDATE RunsData SET
    build_version = $2,
    character_chosen = $3,
    items_purchased_ids = $4,
    items_purged_ids = $5,
    killed_by = $6
WHERE
    id = $1;

-- name: AddFlag :exec
INSERT INTO RunFlags (run_id, flag) VALUES ($1, $2);
-- name: AddCampfire :copyfrom
INSERT INTO CampfireChoice (run_id, "data", floor, "key") VALUES ($1,$2,$3,$4);
-- name: AddDamageTaken :copyfrom
INSERT INTO DamageTaken (run_id, enemies, floor, turns) VALUES ($1,$2,$3,$4);
-- name: AddCardChoice :copyfrom
INSERT INTO CardChoices (run_id, floor, not_picked, picked) VALUES ($1,$2,$3,$4);
-- name: AddRelicObtain :copyfrom
INSERT INTO RelicObtains (run_id, floor, "key") VALUES ($1,$2,$3);
-- name: AddPotionObtain :copyfrom
INSERT INTO PotionObtains (run_id, floor, "key") VALUES ($1,$2,$3);
-- name: AddEventChoices :copyfrom
INSERT INTO EventChoices
    (run_id, damage_delta, event_name_id, floor, gold_delta, max_hp_delta, player_choice_id,
     relics_obtained_ids)
VALUES
    ($1,$2,$3,$4,$5,$6,$7,$8);
-- name: AddMasterDeck :copyfrom
INSERT INTO MasterDecks (run_id, card_id, count, upgrades) VALUES ($1,$2,$3,$4);
-- name: AddBossRelics :copyfrom
INSERT INTO BossRelics (run_id, not_picked, picked) VALUES ($1,$2,$3);

-- name: GetStr :one
SELECT id FROM StrCache WHERE str = $1;

-- name: StrCacheAdd :exec
SELECT str_cache_add($1::text[]);
-- name: StrCacheToId :many
SELECT str_cache_to_id($1::text[]);

-- name: GetCampfires :many
SELECT CC.id, CC.data, CC.floor, StrCache.str as "key" FROM CampfireChoice CC
    LEFT JOIN StrCache ON CC.key = StrCache.id
    WHERE CC.id = $1
    ORDER BY floor;




