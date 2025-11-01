-- +goose Up
CREATE TABLE IF NOT EXISTS job_offers (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    company TEXT,
    location TEXT,
    salary TEXT,
    description TEXT,
    url TEXT UNIQUE NOT NULL,
    source TEXT NOT NULL,
    published_at DATETIME,
    skills TEXT, -- JSON array
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_seen_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS job_offers;
