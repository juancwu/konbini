-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS magic_links (
    token TEXT NOT NULL,
    user_id TEXT NOT NULL,
    created_at TEXT NOT NULL,
    expires_at TEXT NOT NULL,

    CONSTRAINT pk_user_magic_link PRIMARY KEY (user_id, token),
    CONSTRAINT unique_token UNIQUE (token),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS magic_links;
-- +goose StatementEnd