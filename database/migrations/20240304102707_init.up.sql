CREATE TABLE users
(
    id         INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
    name       TEXT    NOT NULL,
    admin      BOOLEAN NOT NULL,
    token_hash TEXT    NOT NULL
);
