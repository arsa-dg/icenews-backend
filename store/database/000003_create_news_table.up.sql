CREATE TYPE scope AS ENUM('top_news', 'breaking_news');

CREATE TABLE news(
  id serial primary key,
  title varchar(255) not null,
  slug_url varchar(255) not null,
  cover_image varchar(255) not null,
  nsfw boolean not null,
  content text not null,
  upvote int not null default 0,
  downvote int not null default 0,
  comment int not null default 0,
  view int not null default 0,
  scope scope not null,
  created_at timestamp not null,
  author_id uuid not null,
  category_id integer not null
);

ALTER TABLE news ADD FOREIGN KEY (author_id) REFERENCES users(id);

ALTER TABLE news ADD FOREIGN KEY (category_id) REFERENCES categories(id);