WITH
  query AS (
    /* @members: []string */
    SELECT [ %[1]s ] AS members
  ),
  candidateGames AS (
    SELECT
      match.details.gameId,
      participant.teamId,
      member
    FROM
      /* <project>.<dataset>.<table> */
      `%[2]s.%[3]s.%[4]s` as match,
      query,
      UNNEST(match.details.participantIdentities) AS identity
    JOIN
      UNNEST(members) AS member
    ON (
      identity.player.currentAccountId = member
      OR identity.player.accountId = member
    )
    JOIN
      UNNEST(match.details.participants) AS participant
    ON
      participant.participantId = identity.participantId
  )

SELECT
  match.details.*
FROM
  /* <project>.<dataset>.<table> */
  `%[2]s.%[3]s.%[4]s` as match
WHERE (
  (
    SELECT COUNT(*)
    FROM candidateGames
  ) > 4
  AND
  (
    SELECT COUNT(DISTINCT teamId)
    FROM candidateGames
    WHERE match.details.gameId = candidateGames.gameId
  ) = 1
)
