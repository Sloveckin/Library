create table Books (
    Id serial not null primary key ,
    Name varchar(64) not null
);

create table Authors (
    Id serial not null primary key,
    BookId integer,
    Name varchar(64) not null,

    foreign key (BookId) references Books(Id)
);
