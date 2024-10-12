CREATE TABLE Users 
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
    title VARCHAR(100) NOT NULL,
    description VARCHAR(100) NOT NULL,
    markdown_url VARCHAR(100) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);

CREATE TABLE Comments (
    id BIGINT NOT NULL,
    post_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    content VARCHAR(500) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES Posts(id) ON DELETE CASCADE
);

INSERT INTO Users (id, username) VALUES (1796290045997481984, "johndoe");
INSERT INTO Posts (id, user_id, created_at, title, description, markdown_url) VALUES (1796290045997481985, 1796290045997481984, "2024-04-04 00:00:00", "My first post", "First post description", "https://example.com/first-post");
INSERT INTO Comments (id, post_id, user_id, content, created_at) VALUES (1796290045997481986, 1796290045997481985, 1796290045997481984, "Example Comment", "2024-04-04 00:00:00");
INSERT INTO Comments (id, post_id, user_id, content, created_at) VALUES (1796290045997481987, 1796290045997481985, 1796290045997481984, "Example Comment 2", "2024-04-05 00:00:00");
INSERT INTO Comments (id, post_id, user_id, content, created_at) VALUES (1796290045997481988, 1796290045997481985, 1796290045997481984, "Example Comment 3", "2024-04-05 20:00:00");
INSERT INTO Comments (id, post_id, user_id, content, created_at) VALUES (1796290045997481989, 1796290045997481985, 1796290045997481984, "Example Comment 4", "2024-04-06 00:00:00");
