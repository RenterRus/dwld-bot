-- +goose Up
-- +goose StatementBegin
create table if not exists links (
link text unique, 
quality integer, 
sendingAt text, 
userID text,
userName text,
messageID text,
errorMsg text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table links;
-- +goose StatementEnd
