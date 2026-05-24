ALTER TABLE users
ADD COLUMN IF NOT EXISTS email TEXT,
ADD COLUMN IF NOT EXISTS password_hash TEXT,
ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT 'user';

UPDATE users
SET email = 'user' || id || '@example.local'
WHERE email IS NULL OR email = '';

UPDATE users
SET password_hash = ''
WHERE password_hash IS NULL;

ALTER TABLE users
ALTER COLUMN email SET NOT NULL,
ALTER COLUMN password_hash SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_active
ON users (email)
WHERE deleted_at IS NULL;
