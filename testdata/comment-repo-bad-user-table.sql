CREATE TABLE Unknown 
(
    id BIGINT NOT NULL,
    username VARCHAR(30) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE Posts
(
    id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    tagline VARCHAR(100) NOT NULL,
    markdown_url VARCHAR(100) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES Unknown(id) ON DELETE CASCADE
);

CREATE TABLE Comments (
    id BIGINT NOT NULL,
    post_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    content VARCHAR(500) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES Unknown(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES Posts(id) ON DELETE CASCADE
);

INSERT INTO Unknown (id, username) VALUES (1796290045997481984, 'johndoe');

INSERT INTO Posts (id, user_id, created_at, tagline, markdown_url) VALUES (1796301682498338816, 1796290045997481984, "2024-04-04 00:00:00", "Example TagLine", "https://example.com/exampleid");

INSERT INTO Comments (id, post_id, user_id, content, created_at) VALUES (1796301682498338817, 1796301682498338816, 1796290045997481984, "Example Comment", "2024-04-04 00:00:00");
