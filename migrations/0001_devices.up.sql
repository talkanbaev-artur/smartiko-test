create table if not exists devices (
    id serial primary key,
    eui varchar(40) not null unique
);

create table if not exists params_records (
    id serial primary key,
    device_id integer not null,
    flag smallint not null,
    value boolean not null default 'false',
    change_time timestamptz not null default ('now'),
    constraint unq_devid_flag UNIQUE (device_id, flag),
	constraint fk_params_records_devices
		foreign key (device_id)
			references devices(id)
);

