CREATE TABLE thoughts (
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    receiver_id     BIGINT NOT NULL REFERENCES users(id),
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    content_text    TEXT,
    content_url     TEXT,
    media_type      TEXT NOT NULL CHECK (media_type IN ('picture', 'voice memo', 'message', 'doodle', 'music')),
    CHECK (user_id != receiver_id),
    CHECK (content_url IS NOT NULL OR media_type = 'message'),
    CHECK (content_text IS NOT NULL OR media_type = 'message')
);

CREATE INDEX thoughts_by_receiver ON thoughts (receiver_id);
CREATE INDEX thoughts_by_sender ON thoughts (user_id);

