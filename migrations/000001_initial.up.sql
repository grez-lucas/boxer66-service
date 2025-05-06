CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email VARCHAR NOT NULL,
  password BYTEA NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX on users(email);

INSERT INTO users (id, email, password ) VALUES
(1, 'john@example.com', gen_random_bytes(10)),
(2, 'mary@example.com', gen_random_bytes(10)),
(3, 'david@example.com', gen_random_bytes(10));
