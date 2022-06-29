CREATE TABLE users(
  id varchar(255) primary key, 
  username varchar(255) not null unique, 
  password varchar(255) not null, 
  name varchar(255) not null, 
  bio varchar(255) not null, 
  web varchar(255) not null, 
  picture varchar(255) not null
);