create table wallets
(
    id         serial
        constraint wallets_pk
            primary key,
    address    varchar not null,
    name       varchar,
    updated_at int
);

create unique index wallets_address_uindex
    on wallet (address);

create unique index wallets_id_uindex
    on wallet (id);

