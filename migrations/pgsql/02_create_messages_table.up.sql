CREATE TABLE IF NOT EXISTS messages
(
    id         SERIAL PRIMARY KEY,
    user_id    BIGINT    NOT NULL,
    content    TEXT      NOT NULL,
    receiver   TEXT      NOT NULL,
    cost       INT       NOT NULL,
    status     INT       NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

ALTER TABLE messages ADD CONSTRAINT message_content_empty CHECK (content <> '');
ALTER TABLE messages ADD CONSTRAINT message_receiver_empty CHECK (receiver <> '');
ALTER TABLE messages ADD CONSTRAINT fk_users_messages FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE;
CREATE INDEX messages_user_id_idx ON messages USING btree (user_id);
CREATE INDEX messages_receiver_idx ON messages USING btree (receiver);
CREATE INDEX messages_created_at_idx ON messages USING brin (created_at);
