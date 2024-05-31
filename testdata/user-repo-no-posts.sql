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
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);

INSERT INTO Users (id, username) VALUES (1796290045997481984, 'johndoe');
INSERT INTO Users (id, username) VALUES (1796290045997481985, 'janedoe');
