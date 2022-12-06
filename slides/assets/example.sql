SELECT 
  jsonPayload.full_name as repo,
  jsonPayload.stargazers_count as stars,
  jsonPayload.forks_count as forks,
  (jsonPayload.forks_count / jsonPayload.stargazers_count) as ratio,
  jsonPayload.language as lang,
  jsonPayload.size as size,
  PARSE_DATETIME('%Y-%m-%dT%H:%M:%SZ', jsonPayload.starred_at) as starred_at,
  jsonPayload.description as description
FROM `<PROJECT_ID>.github_star.*`
where jsonPayload.language = 'Go'
order by ratio desc
