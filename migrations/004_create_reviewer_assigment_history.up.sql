CREATE TABLE IF NOT EXISTS reviewer_assignment_histories (
     assigment_history_id UUID PRIMARY KEY,
     pr_id                TEXT NOT NULL REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
     user_id              TEXT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
     created_at           TIMESTAMPTZ DEFAULT NOW()
);