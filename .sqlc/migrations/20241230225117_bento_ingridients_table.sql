-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS bento_ingridients (
    id TEXT NOT NULL PRIMARY KEY DEFAULT (gen_random_uuid()),
    bento_id TEXT NOT NULL,
    name TEXT NOT NULL,
    value BLOB NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    CONSTRAINT unique_bento_ingridient_name UNIQUE (bento_id, name),
    CONSTRAINT fk_bento_id FOREIGN KEY (bento_id) REFERENCES bentos(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bento_ingridients;
-- +goose StatementEnd
