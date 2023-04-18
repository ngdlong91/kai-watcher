create table delegators
(
    id            serial
        constraint delegators_pk
            primary key,
    address       varchar not null,
    staked_amount float8 default 0,
    updated_at    int
);

create unique index delegators_id_uindex
    on delegators (id);

create index delegators_staked_amount_index
    on delegators (staked_amount desc);

