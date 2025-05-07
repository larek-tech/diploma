-- +goose Up
-- +goose StatementBegin
create schema chat;
create table chat.chat(
    id text primary key,
    user_id bigint not null,
    title text not null default '',
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);
create index chat_user_id on chat.chat(user_id);

create or replace function chat.update_chat()
    returns trigger
as $$
    begin
        new.updated_at := current_timestamp;
        return new;
    end;
$$
language plpgsql;

create trigger trg_update_chat
    before update on chat.chat
    for each row
execute function chat.update_chat();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop trigger trg_update_chat on chat.chat;
drop function chat.update_chat();
drop schema chat;
-- +goose StatementEnd
