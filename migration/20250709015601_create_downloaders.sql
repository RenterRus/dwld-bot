-- +goose Up
-- +goose StatementBegin
create table if not exists downloaders (
name text unique, 
allowedRootLinks text,
host text, 
port integer
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table downloaders;
-- +goose StatementEnd
