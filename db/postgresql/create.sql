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