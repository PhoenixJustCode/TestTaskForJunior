CREATE TABLE
    books (
        id serial not null unique,
        title varchar(255) not null,
        description varchar(255) not null
    );