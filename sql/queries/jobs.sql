-- name: CreateJobOffer :one
INSERT INTO job_offers (
    id, title, company, location, salary, description, url, source, published_at, skills
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetJobOffer :one
SELECT * FROM job_offers WHERE id = ?;

-- name: GetJobOfferByURL :one
SELECT * FROM job_offers WHERE url = ?;

-- name: UpdateJobOffer :one
UPDATE job_offers 
SET 
    title = ?,
    company = ?,
    location = ?,
    salary = ?,
    description = ?,
    published_at = ?,
    skills = ?,
    last_seen_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpsertJobOffer :one
INSERT INTO job_offers (
    id, title, company, location, salary, description, url, source, published_at, skills
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(url) DO UPDATE SET
    title = excluded.title,
    company = excluded.company,
    location = excluded.location,
    salary = excluded.salary,
    description = excluded.description,
    published_at = excluded.published_at,
    skills = excluded.skills,
    last_seen_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: DeleteJobOffer :exec
DELETE FROM job_offers WHERE id = ?;

-- name: ListJobOffers :many
SELECT * FROM job_offers 
ORDER BY created_at DESC 
LIMIT ? OFFSET ?;

-- name: ListRecentJobOffers :many
SELECT * FROM job_offers 
ORDER BY published_at DESC 
LIMIT ?;

-- name: ListJobOffersBySource :many
SELECT * FROM job_offers 
WHERE source = ?
ORDER BY published_at DESC 
LIMIT ? OFFSET ?;

-- name: ListJobOffersByCompany :many
SELECT * FROM job_offers 
WHERE company = ?
ORDER BY published_at DESC;

-- name: ListJobOffersByLocation :many
SELECT * FROM job_offers 
WHERE location LIKE ?
ORDER BY published_at DESC;

-- name: SearchJobOffers :many
SELECT * FROM job_offers 
WHERE 
    title LIKE ? OR 
    company LIKE ? OR 
    description LIKE ? OR
    skills LIKE ? OR
    location LIKE ? OR
    salary LIKE ?

ORDER BY published_at DESC
LIMIT ? OFFSET ?;
