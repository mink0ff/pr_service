CREATE TABLE teams (
    team_id TEXT PRIMARY KEY,
    team_name TEXT UNIQUE NOT NULL
);

CREATE TABLE users (
    user_id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    team_id TEXT NOT NULL REFERENCES teams(team_id)  ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
    );

CREATE TABLE team_members (
    team_id VARCHAR NOT NULL REFERENCES teams(team_id) ON DELETE CASCADE,
    user_id   VARCHAR NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    PRIMARY KEY (team_id, user_id)
);

CREATE INDEX idx_users_team_id
ON users(team_id);