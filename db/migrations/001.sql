create table blacklist
(
    id serial not null,
    subnet cidr not null
);

create unique index blacklist_id_index
    on blacklist (id);

create unique index blacklist_subnet_index
    on blacklist (subnet);


create table whitelist
(
    id serial not null,
    subnet cidr not null
);

create unique index whitelist_id_index
    on whitelist (id);

create unique index whitelist_subnet_index
    on whitelist (subnet);
