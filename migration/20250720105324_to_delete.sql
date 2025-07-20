-- +goose Up
-- +goose StatementBegin
create table if not exists to_delete (
chatID text,
messageID text,
deleteAt text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table to_delete;
-- +goose StatementEnd