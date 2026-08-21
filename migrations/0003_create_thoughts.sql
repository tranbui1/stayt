CREATE TABLE thoughts (
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    receiver_id     BIGINT NOT NULL REFERENCES users(id),
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    content_text    TEXT,
    content_key     TEXT,
    media_type      TEXT NOT NULL CHECK (media_type IN ('picture', 'voice memo', 'message', 'doodle', 'music')),
    viewed          BOOLEAN NOT NULL DEFAULT FALSE,
    CHECK (user_id != receiver_id),

    CHECK (
        (media_type = 'message' AND content_text IS NOT NULL AND content_key IS NULL)
        OR
        (media_type != 'message' AND content_key IS NOT NULL AND content_text IS NULL)
    )
);

CREATE INDEX thoughts_by_sender ON thoughts (user_id);
CREATE INDEX thoughts_unviewed_by_receiver ON thoughts (receiver_id, viewed);

