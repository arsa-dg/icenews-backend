INSERT INTO users(id, username, password, name, bio, web, picture) VALUES('09106f61-d8a8-456e-932c-82134334dc98', 'tester123', '$2a$10$k2u31pLXIP.2PiPbbmryW.pm2wR8Vxn05JRZzOsvui5Ifgh3EfAti', 'tesname', 'tesbio', 'https://example.com', 'https://img.example.com/profile.jpg');
INSERT INTO users(id, username, password, name, bio, web, picture) VALUES('b347951b-d1c6-4ead-8f97-f23ffb27fa37', 'tes', '$2a$10$PB43Vjqn0rUwjMP6G/CjVedhIDFPWOKYOELDgm7PqUL9y3ApHPc7G', 'tesname2', 'tesbio2', 'https://example2.com', 'https://img.example2.com/profile.jpg');

INSERT INTO categories(name) VALUES('Business');
INSERT INTO categories(name) VALUES('Football');

INSERT INTO news(title, slug_url, cover_image, nsfw, content, scope, author_id, category_id) VALUES('judul 1', 'judul-1', 'https://example.com', true, 'lorem ipsum', 'top_news', 'b347951b-d1c6-4ead-8f97-f23ffb27fa37', 1);
INSERT INTO news(title, slug_url, cover_image, nsfw, content, scope, author_id, category_id) VALUES('judul 2', 'judul-2', 'https://example2.com', true, 'lorem ipsum dolor', 'breaking_news', '09106f61-d8a8-456e-932c-82134334dc98', 2);

INSERT INTO news_images(image, news_id) VALUES('https://example.com', 2);
INSERT INTO news_images(image, news_id) VALUES('https://example2.com', 1);
INSERT INTO news_images(image, news_id) VALUES('https://example3.com', 2);
INSERT INTO news_images(image, news_id) VALUES('https://example4.com', 1);
INSERT INTO news_images(image, news_id) VALUES('https://example5.com', 2);

INSERT INTO comments(description, author_id, news_id) VALUES('lorem ipsum', 'b347951b-d1c6-4ead-8f97-f23ffb27fa37', 1);
INSERT INTO comments(description, author_id, news_id) VALUES('lorem ipsum2', '09106f61-d8a8-456e-932c-82134334dc98', 2);
INSERT INTO comments(description, author_id, news_id) VALUES('lorem ipsum3', 'b347951b-d1c6-4ead-8f97-f23ffb27fa37', 1);
INSERT INTO comments(description, author_id, news_id) VALUES('lorem ipsum4', '09106f61-d8a8-456e-932c-82134334dc98', 2);