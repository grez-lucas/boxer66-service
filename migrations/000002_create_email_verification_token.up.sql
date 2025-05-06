CREATE TABLE IF NOT EXISTS email_verification_tokens (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  email VARCHAR NOT NULL,
  verification_token VARCHAR UNIQUE NOT NULL,
  hashed_password_cache_key VARCHAR NOT NULL,
  token_type VARCHAR DEFAULT 'email_verification',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX ON email_verification_tokens(email);
CREATE INDEX ON email_verification_tokens(verification_token);
CREATE INDEX ON email_verification_tokens(expires_at);
