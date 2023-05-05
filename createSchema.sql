create table Services
(
    username varchar(50) not null,
    service  varchar(50) not null,
    login    varchar(50) null,
    password varchar(50) null,
    primary key (username, service)
);
