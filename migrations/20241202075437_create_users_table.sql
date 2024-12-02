-- +goose Up
-- +goose StatementBegin
create table users(
    id bigserial primary key,
    name varchar not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
