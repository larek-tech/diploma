-- +goose Up
-- +goose StatementBegin
create schema domain;
create table domain.source(
    internal_id bigserial primary key,
    external_id text not null default '',
    user_id bigint not null,
    title text not null,
    content bytea not null,
    type int8 not null,
    update_every_period bigint not null default -1,
    cron_week_day int not null default -1,
    cron_month int not null default -1,
    cron_day int not null default -1,
    cron_hour int not null default -1,
    cron_minute int not null default -1,
    credentials bytea not null default ''::bytea,
    status int4 not null default 1,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);
create index domain_ext_id on domain.source using hash(external_id);

create or replace function domain.update_source()
    returns trigger as
$BODY$
begin
    new.updated_at := current_timestamp;
    return new;
end;
$BODY$
    language plpgsql;

create trigger trg_update_source
    before update on domain.source
    for each row
execute function domain.update_source();

create table domain.source_permitted_roles(
    internal_source_id bigint not null,
    role_id bigint not null,
    created_at timestamp not null default current_timestamp
);

create table domain.source_permitted_users(
    internal_source_id bigint not null,
    user_id bigint not null,
    created_at timestamp not null default current_timestamp
);
create index permitted_user_source on domain.source_permitted_users using hash(user_id);

create or replace function domain.permit_user_source()
    returns trigger as
$BODY$
begin
    insert into domain.source_permitted_users(internal_source_id, user_id)
    values (new.internal_id, new.user_id);
    return new;
end;
$BODY$
    language plpgsql;

create trigger trg_permit_user_source
    after insert on domain.source
    for each row
execute function domain.permit_user_source();

create table domain.domain (
    id bigserial primary key,
    title text not null,
    user_id bigint not null,
    source_ids text[] not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create or replace function domain.update_domain()
    returns trigger as
$BODY$
begin
    new.updated_at := current_timestamp;
    return new;
end;
$BODY$
    language plpgsql;

create trigger trg_update_domain
    before update on domain.domain
    for each row
execute function domain.update_domain();

create table domain.domain_permitted_roles(
    domain_id bigint not null,
    role_id bigint not null,
    created_at timestamp not null default current_timestamp
);

create table domain.domain_permitted_users(
    domain_id bigint not null,
    user_id bigint not null,
    created_at timestamp not null default current_timestamp
);
create index permitted_user_domain on domain.domain_permitted_users using hash(user_id);

create or replace function domain.permit_user_domain()
    returns trigger as
$BODY$
begin
    insert into domain.domain_permitted_users(domain_id, user_id)
    values (new.id, new.user_id);
    return new;
end;
$BODY$
    language plpgsql;

create trigger trg_permit_user_domain
    after insert on domain.domain
    for each row
execute function domain.permit_user_domain();

create function domain.get_permitted_sources(uid bigint, rids bigint[])
    returns table (internal_source_id bigint)
language sql
as $$
    select spu.internal_source_id
    from domain.source_permitted_users spu
        where spu.user_id = uid

    union

    select spr.internal_source_id
    from domain.source_permitted_roles spr
        where spr.role_id = any(rids);
$$;

create function domain.get_permitted_domains(uid bigint, rids bigint[])
    returns table (domain_id bigint)
language sql
as $$
    select dpu.domain_id
    from domain.domain_permitted_users dpu
        where dpu.user_id = uid

    union

    select dpr.domain_id
    from domain.domain_permitted_roles dpr
        where dpr.role_id = any(rids);
$$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop function domain.get_permitted_sources(uid bigint, rids bigint[]);
drop function domain.get_permitted_domains(uid bigint, rids bigint[]);
drop trigger trg_permit_user_domain on domain.domain;
drop function domain.permit_user_domain();
drop table domain.domain_permitted_users;
drop table domain.domain_permitted_roles;
drop trigger trg_update_domain on domain.domain;
drop function domain.update_domain();
drop table domain.domain;
drop trigger trg_permit_user_source on domain.source;
drop function domain.permit_user_source();
drop table domain.source_permitted_users;
drop table domain.source_permitted_roles;
drop trigger trg_update_source on domain.source;
drop function domain.update_source();
drop table domain.source;
drop schema domain;
-- +goose StatementEnd
