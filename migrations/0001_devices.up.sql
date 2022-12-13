create table if not exists devices (
    id serial primary key,
    eui varchar(40) not null unique
);

create table if not exists params_records (
    id serial primary key,
    device_id integer not null,
    flag smallint not null,
    change_time timestamptz not null default ('now')
);
