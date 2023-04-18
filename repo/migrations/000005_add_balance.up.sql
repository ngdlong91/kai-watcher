alter table delegators
    add balance float8 default 0 not null;
alter table delegators
    add total_amount float8 default 0 not null;

