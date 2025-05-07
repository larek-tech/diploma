-- +goose Up
-- +goose StatementBegin
create table domain.scenario (
    id bigserial primary key,
    user_id bigint not null,
    use_multiquery bool not null default false,
    n_queries bigint not null default 0,
    query_model_name text not null default '',
    use_rerank bool not null default false,
    reranker_model_name text not null default '',
    reranker_max_length bigint not null default 0,
    reranker_top_k bigint not null default 0,
    llm_model_name text not null default '',
    temperature decimal not null default 0.0,
    top_k bigint not null default 0,
    top_p decimal not null default 0.0,
    system_prompt text not null default '',
    top_n bigint not null default 0,
    threshold decimal not null default 0.0,
    search_by_query bool default false,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);
create index scenario_user_id on domain.scenario (user_id);

create or replace function domain.update_scenario()
    returns trigger as
$BODY$
begin
    new.updated_at := current_timestamp;
    return new;
end;
$BODY$
    language plpgsql;

create trigger trg_update_scenario
    before update on domain.scenario
    for each row
execute function domain.update_scenario();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop trigger trg_update_scenario on domain.scenario;
drop function domain.update_scenario();
drop table domain.scenario;
-- +goose StatementEnd
