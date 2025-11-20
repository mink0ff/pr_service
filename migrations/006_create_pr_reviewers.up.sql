CREATE TABLE IF NOT EXISTS pr_reviewers (
    pr_id INT NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id),
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (pr_id, user_id)
);
