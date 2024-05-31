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
    tagline VARCHAR(100) NOT NULL,
    markdown_url VARCHAR(100) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES Users(id)
);

INSERT INTO Users (id, username) VALUES (1796290045997481984, 'johndoe');
INSERT INTO Users (id, username) VALUES (1796290045997481985, 'janedoe');

INSERT INTO Posts (id, user_id, created_at, tagline, markdown_url) VALUES  (1796301682498338816, 1796290045997481984, "2024-04-04 00:00:00", "Example TagLine", "https://example.com/exampleid");
INSERT INTO Posts (id, user_id, created_at, tagline, markdown_url) VALUES  (1796301682498338817, 1796290045997481984, "2024-04-05 00:00:00", "Example TagLine 2", "https://example.com/exampleid2");
