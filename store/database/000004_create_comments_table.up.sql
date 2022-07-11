CREATE TABLE comments(
  id serial primary key,
  description text not null,
  created_at varchar(255) not null,
  author_id uuid not null,
  news_id bigint not null
);

ALTER TABLE news ADD FOREIGN KEY (author_id) REFERENCES users(id);

ALTER TABLE comments ADD FOREIGN KEY (news_id) REFERENCES news(id);