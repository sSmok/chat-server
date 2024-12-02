-- +goose Up
-- +goose StatementBegin
create table chats_users(
    id bigserial primary key,
    chat_id bigint references chats(id) on delete cascade,
    user_id bigint references users(id) on delete cascade,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);
alter table chats_users add constraint chat_user_uniq unique(chat_id, user_id);
create index idx_user_id on chats_users(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table chats_users;
-- +goose StatementEnd
