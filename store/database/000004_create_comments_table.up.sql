CREATE TABLE comments(
  id serial primary key,
  description text not null,
  created_at varchar(255) not null,
  news_id bigint not null
);

ALTER TABLE comments ADD FOREIGN KEY (news_id) REFERENCES news (id);