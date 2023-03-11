-- name: GetRun :one
SELECT * FROM RunsData WHERE id = $1 LIMIT 1;

-- name: AddRunRaw :one
INSERT INTO RunsData
    (ascension_level, build_version, campfire_rested, campfire_upgraded, character_id, choose_seed,
     circlet_count, floor_reached, gold, killed_by, local_time, neow_bonus_id, neow_cost_id,
     path_per_floor, path_taken, play_id, player_experience, playtime, purchased_purges, score,
     seed_played, seed_source_timestamp, special_seed, "timestamp", victory, win_rate)
VALUES (
    $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26
)
RETURNING RunsData.id;

-- name: DoesRunExist :one
SELECT count(id)::boolean FROM RunsData R WHERE R.play_id = $1;

-- name: AddFlag :exec
INSERT INTO RunFlags (run_id, flag) VALUES ($1, $2);
-- name: AddCampfire :copyfrom
INSERT INTO CampfireChoice (run_id, "data", floor, "key") VALUES ($1,$2,$3,$4);
-- name: AddDamageTaken :copyfrom
INSERT INTO DamageTaken (run_id, enemies, damage, floor, turns) VALUES ($1,$2,$3,$4,$5);
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
INSERT INTO MasterDecks (run_id, card_id, count) VALUES ($1,$2,$3);
-- name: AddBossRelics :copyfrom
INSERT INTO BossRelics (run_id, not_picked, picked, ord) VALUES ($1,$2,$3,$4);
-- name: AddPerFloor :copyfrom
INSERT INTO PerFloorData (run_id, floor, gold, current_hp, max_hp) VALUES ($1,$2,$3,$4,$5);
-- name: AddRunArrays :copyfrom
INSERT INTO RunArrays (run_id, daily_mods, potions_floor_spawned, potions_floor_usage, relic_ids)
VALUES ($1,$2,$3,$4,$5);
-- name: AddItemsPurchased :copyfrom
INSERT INTO ItemsPurchased (run_id, card_id, floor) VALUES ($1,$2,$3);
-- name: AddItemsPurged :copyfrom
INSERT INTO ItemsPurged (run_id, card_id, floor) VALUES ($1,$2,$3);
-- name: AddRunsExtra :exec
INSERT INTO runs_extra (run_id, extra) VALUES ($1,$2);

-- name: GetStr :one
SELECT id FROM StrCache WHERE str = $1;
-- name: StrCacheAdd :exec
SELECT str_cache_add($1::text[]);
-- name: StrCacheToId :many
SELECT str_cache_to_id($1::text[]);

-- name: RunToJson :one
SELECT r.raw::json, r.path_per_floor::text, r.path_taken::text, r.extra::json
FROM run_to_json((SELECT id FROM runsdata WHERE play_id = $1)) r;

-- name: ArchiveBegin :many
UPDATE rawjsonarchive ra SET status = $1 WHERE status = 0 RETURNING ra.*;
-- name: ArchiveComplete :many
UPDATE rawjsonarchive ra SET status = -1 WHERE status = $1 RETURNING ra.id;
-- name: ArchiveAdd :exec
INSERT INTO RawJsonArchive(bdata, play_id) VALUES ($1, $2);
