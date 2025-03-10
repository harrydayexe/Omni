-- Down Migration: Remove the unique constraint on username
ALTER TABLE posts MODIFY COLUMN description VARCHAR(100) NOT NULL;
