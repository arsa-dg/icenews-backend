CREATE TABLE news_images(
  id serial primary key,
  image varchar(500) not null,
  news_id integer not null
);

ALTER TABLE news_images ADD FOREIGN KEY (news_id) REFERENCES news(id) ON DELETE CASCADE;