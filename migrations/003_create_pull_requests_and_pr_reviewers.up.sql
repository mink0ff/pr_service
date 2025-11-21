CREATE TABLE IF NOT EXISTS pull_requests(
    pull_request_id   TEXT PRIMARY KEY,
    pull_request_name TEXT NOT NULL,
    author_id         TEXT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    status            pr_status NOT NULL DEFAULT 'OPEN',
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    merged_at         TIMESTAMPTZ
);

CREATE TABLE pr_reviewers (
    pull_request_id TEXT NOT NULL REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    reviewer_id     TEXT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,

    PRIMARY KEY (pull_request_id, reviewer_id)
    );

CREATE INDEX idx_pr_reviewers_reviewer_id
ON pr_reviewers (reviewer_id);

CREATE INDEX idx_pull_requests_author_id
ON pull_requests(author_id);

CREATE INDEX idx_pull_requests_status
ON pull_requests(status);
