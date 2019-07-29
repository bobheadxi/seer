SELECT DISTINCT
  gameId
FROM
  /* <project>.<dataset>.<table> */
  `%[1]s.%[2]s.%[3]s`
LIMIT
  500
