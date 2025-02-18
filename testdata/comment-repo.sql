INSERT INTO Users (id, username) VALUES (1796290045997481984, "johndoe");

INSERT INTO Posts (id, user_id, created_at, title, description, markdown_url) VALUES (1796290045997481985, 1796290045997481984, "2024-04-04 00:00:00", "My first post", "First post description", "https://example.com/first-post");

INSERT INTO Comments (id, post_id, user_id, content, created_at) VALUES (1796290045997481986, 1796290045997481985, 1796290045997481984, "Example Comment", "2024-04-04 00:00:00");
INSERT INTO Comments (id, post_id, user_id, content, created_at) VALUES (1796290045997481987, 1796290045997481985, 1796290045997481984, "Example Comment 2", "2024-04-05 00:00:00");
INSERT INTO Comments (id, post_id, user_id, content, created_at) VALUES (1796290045997481988, 1796290045997481985, 1796290045997481984, "Example Comment 3", "2024-04-05 20:00:00");
INSERT INTO Comments (id, post_id, user_id, content, created_at) VALUES (1796290045997481989, 1796290045997481985, 1796290045997481984, "Example Comment 4", "2024-04-06 00:00:00");
