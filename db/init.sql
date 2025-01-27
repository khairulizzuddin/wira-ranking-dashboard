-- Accounts
CREATE TABLE Account (
    acc_id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    encrypted_password TEXT NOT NULL,
    secretkey_2fa TEXT NOT NULL
);

-- Characters (8 per account)
CREATE TABLE Character (
    char_id SERIAL PRIMARY KEY,
    acc_id INT REFERENCES Account(acc_id) ON DELETE CASCADE,
    class_id INT CHECK (class_id BETWEEN 1 AND 8),
    UNIQUE (acc_id, class_id)  -- Enforce 1 character per class per account
);

-- Scores (1 per character)
CREATE TABLE Scores (
    score_id SERIAL PRIMARY KEY,
    char_id INT REFERENCES Character(char_id) ON DELETE CASCADE,
    reward_score INT NOT NULL
);

-- Sessions
CREATE TABLE Session (
    session_id TEXT PRIMARY KEY,
    acc_id INT REFERENCES Account(acc_id) ON DELETE CASCADE,
    expiry_datetime TIMESTAMPTZ NOT NULL
);

-- Indexes for search optimization
CREATE INDEX idx_account_username ON Account(username);
CREATE INDEX idx_character_class ON Character(class_id);
CREATE INDEX idx_scores ON Scores(reward_score DESC);