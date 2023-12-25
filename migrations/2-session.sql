-- +migrate Up
CREATE TABLE sessions
(
    usr_id UUID NOT NULL
        CONSTRAINT sessions_usr_id_pkey
            PRIMARY KEY
                CONSTRAINT sessions_usr_id_fkey
                    REFERENCES usrs (id) ON DELETE SET NULL,
    refresh_token TEXT NOT NULL,
    CONSTRAINT sessions_refresh_token UNIQUE (refresh_token),
    expire_at TIMESTAMP NOT NULL
);

-- +migrate Down
DROP TABLE sessions;