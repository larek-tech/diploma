-- +goose Up
-- +goose StatementBegin
create or replace function domain.add_scenario_to_domain()
    returns trigger as
$$
begin
    update domain.domain
    set scenario_ids = array_append(scenario_ids, new.id),
        updated_at = current_timestamp
    where id = new.domain_id;
end;
$$
    language plpgsql;

create trigger trg_add_scenario_to_domain
    after insert on domain.scenario
    for each row
execute function domain.add_scenario_to_domain();

create or replace function domain.remove_scenario_from_domain()
    returns trigger as
$$
begin
    update domain.domain
    set scenario_ids = array_remove(scenario_ids, old.id),
        updated_at = current_timestamp
    where id = old.domain_id;
end;
$$
    language plpgsql;

create trigger trg_remove_scenario_from_domain
    after delete on domain.scenario
    for each row
execute function domain.remove_scenario_from_domain();

create or replace function domain.cleanup_source_permitted_roles()
    returns trigger as
$$
begin
    delete from domain.source_permitted_roles
    where internal_source_id = old.internal_id;
    return old;
end;
$$
    language plpgsql;

create trigger trg_cleanup_source_permitted_roles
    after delete on domain.source
    for each row
execute function domain.cleanup_source_permitted_roles();

create or replace function domain.cleanup_domain_permitted_roles()
    returns trigger as
$$
begin
    delete from domain.domain_permitted_roles
    where domain_id = old.id;
    return old;
end;
$$
    language plpgsql;

create trigger trg_cleanup_domain_permitted_roles
    after delete on domain.domain
    for each row
execute function domain.cleanup_domain_permitted_roles();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop trigger trg_cleanup_domain_permitted_roles on domain.domain;
drop function domain.cleanup_domain_permitted_roles();
drop trigger trg_cleanup_source_permitted_roles on domain.source;
drop function domain.cleanup_source_permitted_roles();
drop trigger trg_remove_scenario_from_domain on domain.scenario;
drop function domain.remove_scenario_from_domain();
drop trigger trg_add_scenario_to_domain on domain.scenario;
drop function domain.add_scenario_to_domain();
-- +goose StatementEnd
