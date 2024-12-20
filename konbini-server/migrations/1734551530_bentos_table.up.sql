CREATE TABLE bentos (
    id TEXT NOT NULL PRIMARY KEY DEFAULT (uuid4()),
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'utc')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now', 'utc')),

    CONSTRAINT unique_bento_name_user UNIQUE (user_id, name)
);
