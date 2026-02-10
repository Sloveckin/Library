create table books (
    id serial not null primary key,
    name varchar(64) not null
);

create table authors (
    id serial not null primary key,
    bookid integer,
    name varchar(64) not null,
    foreign key (bookid) references books(id)
);
