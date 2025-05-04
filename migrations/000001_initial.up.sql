CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email VARCHAR NOT NULL,
  password VARCHAR NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX on users(email);

INSERT INTO users (id, email, password, created_at, updated_at) VALUES
(1, 'john@example.com', 'example', '2025-05-03 14:23:45Z', '2025-05-03 14:23:45Z'),
(2, 'mary@example.com', 'example', '2025-05-03 14:23:45Z', '2025-05-03 14:23:45Z'),
(3, 'david@example.com', 'example', '2025-05-03 14:23:45Z', '2025-05-03 14:23:45Z');
