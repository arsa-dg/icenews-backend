CREATE TABLE comments(
  id serial primary key,
  description varchar(255) not null,
  created_at timestamp not null default CURRENT_TIMESTAMP,
  author_id uuid not null,
  news_id integer not null
);

ALTER TABLE comments ADD FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE comments ADD FOREIGN KEY (news_id) REFERENCES news(id) ON DELETE CASCADE;