-- Up Migration: Make username field unique
ALTER TABLE posts MODIFY COLUMN markdown_url VARCHAR(255) NOT NULL;
