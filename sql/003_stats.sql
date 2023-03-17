
CREATE VIEW stats_overview AS (
    WITH
        co AS (SELECT '{0.25, 0.5, 0.75}'::float[] as p_quart),
        deck_size AS (
            SELECT R.id, count(D.ix) as total
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
        count(D.ix) as total,
        sum(case when S.upgrades > 0 then 1 else 0 end) as upgrades
    FROM RunsData R
    JOIN MasterDecks D ON R.id = D.run_id
    JOIN CardSpecs S ON D.card_id = S.id
    GROUP BY R.character_id, D.card_id
    ORDER BY D.card_id
);

CREATE FUNCTION per_character_card_stats(char_id int) RETURNS
    TABLE(card_id int, card text, runs int, wins int, deck float4[], floor float4[])
LANGUAGE SQL AS $$
WITH ru AS (SELECT r.id, r.floor_reached, r.victory, count(m.id) as deck_size
            FROM runsdata r
                     INNER JOIN masterdecks m on r.id = m.run_id
            WHERE r.character_id = char_id
            GROUP BY r.id),
     ca AS (SELECT m.card_id, s.card, m.run_id
            FROM masterdecks m
                     LEFT JOIN cardspecsex s on m.card_id = s.id)
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
