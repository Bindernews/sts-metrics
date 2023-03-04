
-- Materialized so we don't have scan RunsData for the list of unique characters.
CREATE MATERIALIZED VIEW IF NOT EXISTS character_list AS (
    SELECT DISTINCT S.id, S.str as "name" FROM StrCache S
    LEFT JOIN RunsData R ON S.id = R.character_id
    WHERE S.id = R.character_id
    ORDER BY S.str
);

CREATE FUNCTION character_list_refresh() RETURNS TRIGGER
LANGUAGE plpgsql AS $$
BEGIN
    REFRESH MATERIALIZED VIEW character_list;
    RETURN NULL;
END $$;

CREATE TRIGGER character_list_refresh AFTER INSERT OR UPDATE OR DELETE
ON RunsData FOR EACH STATEMENT EXECUTE FUNCTION character_list_refresh();

CREATE VIEW stats_overview AS (
    WITH
        co AS (SELECT '{0.25, 0.5, 0.75}'::float[] as p_quart),
        deck_size AS (
            SELECT R.id, sum(D.count) as total
            FROM RunsData R INNER JOIN MasterDecks D ON R.id = D.run_id
            GROUP BY R.id
        )
    SELECT
        CL.id,
        CL.name,
        count(R.id) as runs,
        sum(R.victory::int) as wins,
        avg(R.win_rate) as avg_win_rate,
        percentile_cont(co.p_quart) WITHIN GROUP (ORDER BY deck_size.total)::float4[] as p_deck_size,
        percentile_cont(co.p_quart) WITHIN GROUP (ORDER BY R.floor_reached)::float4[] as p_floor_reached
    FROM character_list CL CROSS JOIN co
    JOIN RunsData R ON R.character_id = CL.id
    JOIN deck_size ON R.id = deck_size.id
    GROUP BY CL.id, CL.name, co.p_quart
);

CREATE VIEW stats_card_counts AS (
    SELECT
        R.character_id as char_id,
        D.card_id,
        sum(D.count) as total,
        sum(case when D.upgrades > 0 THEN D.count else 0 end) as upgrades
    FROM RunsData R
    JOIN MasterDecks D ON R.id = D.run_id
    GROUP BY R.character_id, D.card_id
    ORDER BY D.card_id
);

CREATE FUNCTION per_character_card_stats(char_id int) RETURNS
    TABLE(card_id int, card text, runs int, wins int, deck float4[], floor float4[])
LANGUAGE SQL AS $$
WITH ru AS (SELECT r.id, r.floor_reached, r.victory, sum(m.count) as deck_size
            FROM runsdata r
                     INNER JOIN masterdecks m on r.id = m.run_id
            WHERE r.character_id = char_id
            GROUP BY r.id),
     ca AS (SELECT m.card_id, s.str, m.run_id
            FROM masterdecks m
                     LEFT JOIN strcache s on m.card_id = s.id)
SELECT ca.card_id,
       ca.str               as card,
       count(ru.id)         as runs,
       sum(ru.victory::int) as wins,
       percentile_cont('{0.25, 0.5, 0.75}'::float[]) WITHIN GROUP (ORDER BY ru.deck_size)
           ::float4[]       as deck,
       percentile_cont('{0.25, 0.5, 0.75}'::float[]) WITHIN GROUP (ORDER BY ru.floor_reached)
           ::float4[]       as floor
FROM ca
         JOIN ru ON ru.id = ca.run_id
GROUP BY ca.card_id, ca.str
$$;

---- create above / drop below ----

drop view if exists stats_overview;
drop view if exists stats_card_counts;
drop function if exists per_character_card_stats;
drop trigger if exists character_list_refresh ON RunsData;
drop function if exists character_list_refresh;
drop materialized view if exists character_list;
