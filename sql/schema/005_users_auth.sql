-- +goose Up
ALTER TABLE users
  ADD COLUMN email TEXT UNIQUE NOT NULL,
  ADD COLUMN password TEXT NOT NULL;

CREATE TABLE tokens (
    user_id UUID NOT NULL,
    CONSTRAINT fk_user_id
      FOREIGN KEY (user_id)
      REFERENCES users(id)
      ON DELETE CASCADE,
    value TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL
);

-- +goose Down
ALTER TABLE users
  DROP COLUMN email,
  DROP COLUMN password;

DROP TABLE tokens;