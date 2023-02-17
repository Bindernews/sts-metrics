

-- name: StatsGetCharRuns :many
SELECT R.id FROM RunsData R
LEFT JOIN strcache s on R.character_chosen = s.id
WHERE s.str = $1;

-- name: StatsGetOverall :many
WITH
    c_runs AS (SELECT * FROM RunsData WHERE character_chosen = $1)
SELECT 'runs', count(id) FROM c_runs
UNION SELECT 'avg_win_rate', avg(win_rate) FROM c_runs
UNION SELECT 'avg_cards', avg(array_length(master_deck, 1)) FROM c_runs;

