-- +goose Up
-- +goose StatementBegin
create schema chat;
create table chat.chat(
    id uuid primary key default gen_random_uuid(),
    user_id bigint not null,
    title text not null default '',
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    is_deleted bool not null default false
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

create table chat.query(
    id bigserial primary key,
    user_id bigint not null,
    chat_id uuid not null,
    content text not null,
    domain_id bigint not null default 0,
    source_ids text[] not null default array[]::text[],
    scenario_id bigint not null default 0,
    metadata jsonb default '{}'::jsonb,
    created_at timestamp not null default current_timestamp
);

create table chat.response(
    id bigserial primary key,
    query_id bigint not null,
    chat_id uuid not null,
    content text not null,
    status int8 not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create or replace function chat.update_response()
    returns trigger
as $$
begin
    new.updated_at := current_timestamp;
    return new;
end;
$$
    language plpgsql;

create trigger trg_update_response
    before update on chat.response
    for each row
execute function chat.update_response();

create or replace function chat.touch_chat_on_query()
    returns trigger as $$
begin
    update chat.chat
    set updated_at = current_timestamp
    where id = new.chat_id;
    return new;
end;
$$ language plpgsql;

create trigger trg_update_chat_by_query
    after insert on chat.query
    for each row
execute function chat.touch_chat_on_query();

create or replace function chat.touch_chat_on_response()
    returns trigger as $$
begin
    update chat.chat
    set updated_at = current_timestamp
    where id = new.chat_id;
    return new;
end;
$$ language plpgsql;

create trigger trg_update_chat_by_response
    after insert on chat.response
    for each row
execute function chat.touch_chat_on_response();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop trigger trg_update_chat_by_query on chat.query;
drop function chat.touch_chat_on_response();
drop trigger trg_update_chat_by_response on chat.response;
drop function chat.touch_chat_on_query();
drop trigger trg_update_response on chat.response;
drop function chat.update_response();
drop table chat.response;
drop table chat.query;
drop trigger trg_update_chat on chat.chat;
drop function chat.update_chat();
drop table chat.chat;
drop schema chat;
-- +goose StatementEnd
