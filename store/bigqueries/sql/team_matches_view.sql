/*
  team_matches_view.sql

  This query generates a view where each row is contains data for a player in
  a single match, such that each game maps to several rows in this view.
*/

WITH
  query AS (
    /* @members: []string */
    SELECT [ %[1]s ] AS members
  ),
  /*
    candidateGames is a table of games that have members of the query team as
    participants
  */
  candidateGames AS (
    SELECT
      match.gameId,
      member,
      participant
    FROM
      /* <project>.<dataset>.<table> */
      `%[2]s.%[3]s.%[4]s` as match,
      query,
      UNNEST(match.participantIdentities) AS identity
    JOIN
      UNNEST(members) AS member
    ON (
      identity.player.currentAccountId = member
      OR identity.player.accountId = member
    )
    JOIN
      UNNEST(match.participants) AS participant
    ON
      participant.participantId = identity.participantId
  )

SELECT
  /* identifying information */
  candidateGames.gameId,
  candidateGames.member AS accountId,
  /* match metadata */
  match.seasonId,
  match.queueId,
  match.gameDuration,
  match.gameCreation,
  match.gameVersion,
  /* participant data */
  candidateGames.participant.*,
  /* frames for member in match */
  ARRAY(
    SELECT AS STRUCT
      frame.timestamp timestamp,
      frame.participantFrames[ORDINAL(candidateGames.participant.participantId)] participantFrames
    FROM
      `seer-engine.dev.timelines_integration_test` AS timeline,
      UNNEST(timeline.frames) AS frame
    WHERE
      timeline.gameId = candidateGames.gameId
  ) AS frames
FROM
  /* <project>.<dataset>.<table> */
  `%[2]s.%[3]s.%[4]s` as match
LEFT JOIN
  candidateGames
ON
  match.gameId = candidateGames.gameId
WHERE (
  (
    /* 4 of more members from the team are in the game */
    SELECT COUNT(*)
    FROM candidateGames
  ) > 4
  AND
  (
    /* all members are on the same team */
    SELECT COUNT(DISTINCT participant.teamId)
    FROM candidateGames
    WHERE match.gameId = candidateGames.gameId
  ) = 1
)
ORDER BY gameCreation DESC
