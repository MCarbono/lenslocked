-- +goose Up
-- +goose StatementBegin
CREATE TABLE galleries(
    id TEXT PRIMARY KEY,
    user_id TEXT UNIQUE REFERENCES users (id) ON DELETE CASCADE,
    title TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE galleries;
-- +goose StatementEnd
