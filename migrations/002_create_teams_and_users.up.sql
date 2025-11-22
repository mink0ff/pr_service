CREATE TABLE IF NOT EXISTS teams (
    team_id TEXT PRIMARY KEY,
    team_name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    user_id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    team_id TEXT NOT NULL REFERENCES teams(team_id)  ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
    );

CREATE INDEX idx_users_team_id
ON users(team_id);