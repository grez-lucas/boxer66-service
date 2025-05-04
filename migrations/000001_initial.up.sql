CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email VARCHAR NOT NULL,
  password BYTEA NOT NULL,
  salt BYTEA,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX on users(email);

INSERT INTO users (id, email, password, created_at, updated_at) VALUES
(1, 'john@example.com', '\x$2a$10$someRandomSaltedHashForJohn', now(), now()),
(2, 'mary@example.com', '\x$2a$10$anotherDifferentSaltedHashForMary', now(), now()),
(3, 'david@example.com', '\x$2a$10$anotherDifferentSaltedHashForDavid', now(), now());
