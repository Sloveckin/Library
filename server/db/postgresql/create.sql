-- create table books (
--     id serial not null primary key,
--     name varchar(64) not null
-- );

-- create table authors (
--     id serial not null primary key,
--     bookid integer,
--     name varchar(64) not null,
--     foreign key (bookid) references books(id)
-- );

create table Books (
    Id serial not null primary key ,
    Name varchar(64) not null
);

create table Authors (
    Id serial not null primary key,
    Name varchar(64) not null
);

create table AuthorToBook(
    AuthorId integer,
    BookId integer,

    foreign key (AuthorId) references Authors(Id),
    foreign key (BookId) references Books(Id)
);