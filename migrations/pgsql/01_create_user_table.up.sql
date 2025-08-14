CREATE TABLE IF NOT EXISTS users
(
    id         SERIAL PRIMARY KEY,
    name       varchar(250) NOT NULL,
    balance    BIGINT       NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

ALTER TABLE users ADD CONSTRAINT user_name_empty CHECK (name <> '');
ALTER TABLE users ADD CONSTRAINT user_insufficient_balance CHECK (balance >= 0);


