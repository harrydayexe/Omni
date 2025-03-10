-- Up Migration: Make username field unique
ALTER TABLE posts MODIFY COLUMN description VARCHAR(255) NOT NULL;
