-- +goose Up
-- +goose StatementBegin
create table messages(
    id bigserial primary key,
    source_id bigint references chats_users(id) on delete cascade,
    text text not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table messages;
-- +goose StatementEnd
