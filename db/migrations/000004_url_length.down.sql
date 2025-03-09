-- Down Migration: Remove the unique constraint on username
ALTER TABLE posts MODIFY COLUMN markdown_url VARCHAR(100) NOT NULL;
