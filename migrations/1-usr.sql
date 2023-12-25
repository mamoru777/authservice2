-- +migrate Up
CREATE TABLE usrs
(
    id UUID DEFAULT UUID_GENERATE_V4()
        CONSTRAINT usrs_usr_id_pkey
            PRIMARY KEY,
    login TEXT NOT NULL,
    CONSTRAINT usrs_login UNIQUE (login),
    email TEXT NOT NULL,
    CONSTRAINT usrs_email UNIQUE (email),
    password BYTEA NOT NULL,
    isSignedUp BOOLEAN NOT NULL
);

-- +migrate Down
DROP TABLE usrs;
