-- +goose Up
-- +goose StatementBegin
alter database dev
    set timezone to 'Europe/Moscow';
create schema auth;
create table auth.role (
    id bigserial primary key,
    name text unique not null,
    created_at timestamp not null default current_timestamp,
    is_deleted bool not null default false
);
create index role_name on auth.role using hash(name);

insert into auth.role(name)
values ('default'),
       ('admin');

create table auth.user (
    id bigserial primary key,
    email text unique not null,
    hash_password text not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    is_deleted bool not null default false
);
create index user_email on auth.user using hash(email);

create or replace function auth.update_user()
    returns trigger as
$BODY$
begin
    new.updated_at := current_timestamp;
    return new;
end;
$BODY$
    language plpgsql;

create trigger trg_update_user
    before update on auth.user
    for each row
execute function auth.update_user();

create table auth.user_role (
    user_id bigint not null,
    role_id bigint not null,
    created_at timestamp not null default current_timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter database dev
    set timezone to 'UTC';
drop trigger trg_update_user on auth.user;
drop function auth.update_user;
drop table auth.user_role;
drop table auth.user;
drop table auth.role;
drop schema auth;
-- +goose StatementEnd
