
-- name: StatsGetOverview :one
SELECT * FROM stats_overview WHERE "name" = $1;

-- name: StatsListCharacters :many
SELECT * FROM character_list;
