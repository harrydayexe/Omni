-- Up Migration: Make username field unique
ALTER TABLE users ADD CONSTRAINT unique_username UNIQUE (username);
CREATE INDEX idx_username ON users(username);
