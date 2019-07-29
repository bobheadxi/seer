WITH
  query AS (
    /* @members: []string */
    SELECT [ %[1]s ] AS members
  ),
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
  match.gameId, candidateGames.member AS accountId, 
  match.seasonId, match.queueId, match.gameDuration, match.gameCreation, match.gameVersion,
  candidateGames.participant.*,
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
  `%[2]s.%[3]s.%[4]s` as match,
  candidateGames
WHERE (
  (
    SELECT COUNT(*)
    FROM candidateGames
  ) > 4
  AND
  (
    SELECT COUNT(DISTINCT participant.teamId)
    FROM candidateGames
    WHERE match.gameId = candidateGames.gameId
  ) = 1
)
