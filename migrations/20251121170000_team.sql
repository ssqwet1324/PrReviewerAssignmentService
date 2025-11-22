-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS team (
    team_name TEXT PRIMARY KEY
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS team;
-- +goose StatementEnd
