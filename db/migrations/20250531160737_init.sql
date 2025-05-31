-- +goose Up
-- +goose StatementBegin
CREATE TABLE authors (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    bio TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE authors;
-- +goose StatementEnd
