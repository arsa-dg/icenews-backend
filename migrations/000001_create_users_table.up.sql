CREATE TABLE users(
  id uuid primary key, 
  username varchar(20) not null unique, 
  password varchar(255) not null, 
  name varchar(255) not null, 
  bio varchar(255) not null, 
  web varchar(255) not null, 
  picture varchar(255) not null
);