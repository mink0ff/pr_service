CREATE TABLE IF NOT EXISTS pull_requests (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author_id INT NOT NULL REFERENCES users(id),
    team_id INT NOT NULL REFERENCES teams(id),
    status pr_status NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);