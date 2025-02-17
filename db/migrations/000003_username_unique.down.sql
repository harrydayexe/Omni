-- Down Migration: Remove the unique constraint on username
ALTER TABLE users DROP CONSTRAINT unique_username;
DROP INDEX idx_username;
