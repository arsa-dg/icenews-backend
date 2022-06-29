CREATE TABLE news(
  id serial primary key,
  title varchar(255) not null,
  slug_url varchar(255) not null,
  cover_image varchar(255) not null,
  nsfw boolean not null,
  content text not null,
  created_at varchar(255) not null,
  author_id varchar(255),
  category_id bigint
);

ALTER TABLE news ADD FOREIGN KEY (author_id) REFERENCES users(id);

ALTER TABLE news ADD FOREIGN KEY (category_id) REFERENCES categories(id);