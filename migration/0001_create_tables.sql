create table if not exists accounts (
    id         serial primary key,
    balance    int       not null default 0,
    created_at timestamp not null default now(),
    updated_at timestamp     default null,
    deleted_at timestamp     default null
);

create table if not exists products (
    id         serial primary key,
    name       varchar(255) not null unique,
    created_at timestamp not null default now(),
    updated_at timestamp     default null,
    deleted_at timestamp     default null
);

create table if not exists reservations (
    id         serial primary key,
    account_id int       not null,
    product_id int       not null,
    amount     int       not null,
    created_at timestamp not null default now(),
    updated_at timestamp     default null,
    deleted_at timestamp     default null,
    foreign key (account_id) references accounts (id),
    foreign key (product_id) references products (id)
);

create table if not exists operations (
    id             serial primary key,
    account_id     int          not null,
    amount         int          not null,
    operation_type varchar(255) not null,
    product_id     int                   default null,
    description    varchar(255)          default null,
    created_at     timestamp not null default now(),
    updated_at     timestamp     default null,
    deleted_at     timestamp     default null,
    foreign key (account_id) references accounts (id)
);
