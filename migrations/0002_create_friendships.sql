CREATE TABLE friendships (
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id),
    friend_id       BIGINT NOT NULL REFERENCES users(id),
    requested_by    BIGINT NOT NULL REFERENCES users(id), 
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    status          TEXT NOT NULL CHECK (status IN ('pending', 'rejected', 'accepted')),
    CHECK (requested in (user_id, friend_id)),
    CHECK (user_id < friend_id),
    UNIQUE (user_id, friend_id)
);