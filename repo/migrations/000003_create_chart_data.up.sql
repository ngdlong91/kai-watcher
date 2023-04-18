create table chart_data
(
    id                   integer generated by default as identity,
    chart_id             integer                 not null,
    date                 date                    not null,
    value                varchar(64)             not null,
    created_at           timestamp default now() not null,
    min_blockscout_block bigint,
    primary key (id),
    foreign key (chart_id) references charts
);

comment on table chart_data is 'Table contains chart data points';

alter table chart_data
    owner to postgres;

create unique index chart_data_chart_id_date_idx
    on chart_data (chart_id, date);

