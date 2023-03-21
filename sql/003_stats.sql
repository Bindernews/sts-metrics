
CREATE VIEW stats_overview AS (
    WITH
        co AS (SELECT '{0.25, 0.5, 0.75}'::float[] as p_quart),
        deck_size AS (
            SELECT R.id, array_length(a.master_deck, 1) as total
            FROM RunsData R INNER JOIN RunArrays a ON R.id = a.run_id
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
    WITH deck AS (
        SELECT r.character_id as char_id,
               unnest(a.master_deck) as card_id
        FROM RunsData r
        JOIN RunArrays a ON r.id = a.run_id
    )
    SELECT d.char_id,
           d.card_id,
           count(d.card_id) as total,
           sum(case when s.upgrades > 0 then 1 else 0 end) as upgrades
    FROM deck d
    JOIN CardSpecs s ON d.card_id = s.id
    GROUP BY d.char_id, D.card_id
    ORDER BY d.card_id
);

CREATE FUNCTION per_character_card_stats(char_id int) RETURNS
    TABLE(card_id int, card text, runs int, wins int, deck float4[], floor float4[])
LANGUAGE SQL AS $$
WITH ru AS (SELECT r.id,
                   r.floor_reached,
                   r.victory,
                   array_length(a.master_deck, 1) as deck_size
            FROM runsdata r
                     INNER JOIN runarrays a on r.id = a.run_id
            WHERE r.character_id = char_id),
     ca AS (SELECT s.id as card_id,
                   s.card,
                   a.run_id
            FROM runsdata r
                    INNER JOIN runarrays a on a.run_id = r.id
                    LEFT JOIN cardspecsex s on s.id = any(a.master_deck))
SELECT ca.card_id,
       ca.card,
       count(ru.id)         as runs,
       sum(ru.victory::int) as wins,
       percentile_cont('{0.25, 0.5, 0.75}'::float[]) WITHIN GROUP (ORDER BY ru.deck_size)
           ::float4[]       as deck,
       percentile_cont('{0.25, 0.5, 0.75}'::float[]) WITHIN GROUP (ORDER BY ru.floor_reached)
           ::float4[]       as floor
FROM ca
         JOIN ru ON ru.id = ca.run_id
GROUP BY ca.card_id, ca.card
$$;

CREATE FUNCTION card_pick_stats(char_id int, merge_upgrades bool) RETURNS
    TABLE(card text, pick int, skip int)
LANGUAGE SQL IMMUTABLE AS $$
with cc as (select c.id, c.not_picked, c.picked
            from cardchoices c
                     inner join runsdata r on r.id = c.run_id
            where r.character_id = char_id),
     -- Expand and count not-picked cards
     c_not as (select np.id, count(np.id) n
               from (select unnest(cc.not_picked) id from cc) np
               group by np.id),
     -- Expand and count picked cards
     c_pick as (select cc.picked as id, count(cc.picked) n
                from cc
                group by cc.picked),
     -- select card name based on merge_upgrades value
     cf as (select cn.id,
                   coalesce(cp.n, 0)     as pick,
                   cn.n                  as skip,
                   (case
                        when merge_upgrades then s.card
                        else s.card_full end) as card
            from c_not cn
                     full join c_pick cp on cn.id = cp.id
                     join cardspecsex s on s.id = cn.id)
-- re-sum, grouping by card name
select cf.card, sum(cf.skip), sum(cf.pick)
from cf
group by card
$$;

---- create above / drop below ----

drop view if exists stats_overview;
drop view if exists stats_card_counts;
drop function if exists per_character_card_stats;
drop function if exists card_pick_stats;
