CREATE TABLE artifacts
(
    mvn_group TEXT     NOT NULL,
    artifact  TEXT     NOT NULL,
    version   TEXT     NOT NULL,
    modified  DATETIME NOT NULL
);

CREATE TABLE paths
(
    user_id    INTEGER UNIQUE,
    path       TEXT PRIMARY KEY,
    deploy     TINYINT,
    created_at DATETIME,
    updated_at DATETIME,

    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE users
(
    id         INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
    name       TEXT    NOT NULL,
    admin      BOOLEAN NOT NULL,
    token_hash TEXT    NOT NULL
);
