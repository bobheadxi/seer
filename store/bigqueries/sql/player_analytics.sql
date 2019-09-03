/* TODO: templatize */

CREATE TEMP FUNCTION COUNTS(vals ARRAY<STRING>)
RETURNS ARRAY<STRUCT<value STRING, count INT64>>
AS ((
  SELECT ARRAY(
    SELECT AS STRUCT 
      vs AS value,
      COUNT(vs) AS count
    FROM UNNEST(vals) AS vs
    GROUP BY vs
    ORDER BY count DESC
  )
));

WITH
  peakTiers AS (
    SELECT accountId, gameCreation, seasonTier FROM (
      SELECT DISTINCT
        accountId,
        gameCreation,
        highestAchievedSeasonTier AS seasonTier,
        RANK() OVER ( PARTITION BY accountId ORDER BY gameCreation ) AS _fn_rank
      FROM
        /* <project>.<dataset>.<table> */
        `seer-engine.dev.team_test_team`
      WHERE highestAchievedSeasonTier IS NOT NULL
    )
    WHERE _fn_rank = 1
  ),
  playerAggregations AS (
    SELECT
      accountId,
      ARRAY_AGG(STRUCT(
        championId,
        count,
        lanes,
        roles,
        perks,
        summoners
      ) ORDER BY count DESC) AS champions
    FROM (
      SELECT
        t.accountId,
        t.championId,
        COUNT(*) AS count,
        COUNTS(ARRAY_AGG(timeline.lane)) AS lanes,
        COUNTS(ARRAY_AGG(timeline.role)) AS roles,
        COUNTS(ARRAY_AGG(FORMAT("%d,%d", stats.perkPrimaryStyle, stats.perkSubStyle))) AS perks,
        COUNTS(ARRAY_AGG(FORMAT("%d,%d", spell1Id, spell2Id))) AS summoners
      FROM   
        /* <project>.<dataset>.<table> */
        `seer-engine.dev.team_test_team` as t
      GROUP BY accountId, championId
    )
    GROUP BY accountId
    ORDER BY accountId
  )

SELECT
  playerAggregations.accountId,
  peakTiers.seasonTier,
  playerAggregations
FROM playerAggregations
LEFT JOIN peakTiers ON playerAggregations.accountId = peakTiers.accountId
