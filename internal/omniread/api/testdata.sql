INSERT INTO users 
    (id, username, password) 
VALUES 
    (1796290045997481984, "johndoe", "$2a$10$dPMvZy6/IXgBSb1KVB.rD.fRy5V3OIcrg1GYzKDCT2Motcth./gV6"),
    (1796290045997481985, "janedoe", "$2a$10$gQqRIqRSut1YBNQOz0Sd5eAdATDLHqkhGdUDTbLWzzP.mmA7Fmo5G"),
    (1796290045997481986, "harrydayexe", "$2a$10$CXVhB7idJ6wfOorTyCl9L.kazwTam.g8ai1Iq1..0CFd/rXypFiG."),
    (1796290045997481987, "johnsmith", "$2a$10$E0hoUrdlphfpCMBSFKjeR.1ILtuK8Gi4erefXDXW4FeAdN/4Tctsa"),
    (1796290045997481988, "testername", "$2a$10$NmlDbJyVFMUxDtqOcZaXuunknlrhb7insV8mwsZEtD.Po/VLrun7a");

INSERT INTO posts 
    (id, user_id, created_at, title, description, markdown_url) 
VALUES 
    (1796290045997481995, 1796290045997481984, "2024-04-04 00:00:00", "My first post", "First post description", "https://example.com/johndoe-first-post"),
    (1796290045997481996, 1796290045997481984, "2024-05-04 00:00:00", "My second post", "Second post description", "https://example.com/johndoe-second-post"),
    (1796290045997481997, 1796290045997481985, "2024-06-04 00:00:00", "My first post", "First post description", "https://example.com/janedoe-first-post"),
    (1796290045997481998, 1796290045997481985, "2024-06-04 00:00:00", "My second post", "First post description", "https://example.com/janedoe-first-post"),
    (1796290045997481999, 1796290045997481985, "2024-06-04 00:00:00", "My third post", "First post description", "https://example.com/janedoe-first-post"),
    (1796290045997482990, 1796290045997481985, "2024-07-04 00:00:00", "My fourth post", "First post description", "https://example.com/janedoe-first-post"),
    (1796290045997482991, 1796290045997481985, "2024-08-04 00:00:00", "My fifth post", "First post description", "https://example.com/janedoe-first-post"),
    (1796290045997482992, 1796290045997481985, "2024-09-04 00:00:00", "My sixth post", "First post description", "https://example.com/janedoe-first-post"),
    (1796290045997482993, 1796290045997481985, "2024-10-04 00:00:00", "My seventh post", "First post description", "https://example.com/janedoe-first-post"),
    (1796290045997482994, 1796290045997481985, "2024-11-04 00:00:00", "My eighth post", "First post description", "https://example.com/janedoe-first-post"),
    (1796290045997482995, 1796290045997481985, "2024-12-04 00:00:00", "My ninth post", "First post description", "https://example.com/janedoe-first-post"),
    (1796290045997482996, 1796290045997481985, "2024-12-05 00:00:00", "My tenth post", "First post description", "https://example.com/janedoe-first-post"),
    (1796290045997482997, 1796290045997481985, "2024-12-06 00:00:00", "My eleventh post", "First post description", "https://example.com/janedoe-first-post"),
    (1796290045997482998, 1796290045997481986, "2024-12-07 00:00:00", "My twelfth post", "First post description", "https://example.com/harrydayexe-first-post");

INSERT INTO comments 
    (id, post_id, user_id, content, created_at) 
VALUES 
    (1796290045997481886, 1796290045997481995, 1796290045997481986, "Example Comment", "2024-04-04 00:00:00"),
    (1796290045997481887, 1796290045997481996, 1796290045997481987, "Example Comment 2", "2024-04-05 00:00:00"),
    (1796290045997481888, 1796290045997481997, 1796290045997481987, "Example Comment 3", "2024-04-05 20:00:00"),
    (1796290045997481889, 1796290045997481997, 1796290045997481988, "Example Comment 4", "2024-04-06 00:00:00");
