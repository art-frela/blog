drop database if exists blog;
create database blog;
-- create table scope
use blog;
drop table if exists users;
create table users
(
    id         varchar(42) PRIMARY KEY,
    username       varchar(255)                       null,
    nick       varchar(255)                       null,
    email   varchar(500)                               null,
    created_at datetime default CURRENT_TIMESTAMP null,
    modified_at datetime default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    user_role int default -1 not null,
    salt varchar(25) default 'saltsalt' not null,
    avatar_url varchar(512) null
);

-- insert default user
insert into users (id, username, nick, email, avatar_url) VALUES ('00000000-0000-0000-00000000', 'anonimous', 'anonimous', 'user@example.com', 'https://getuikit.com/docs/images/avatar.jpg');

drop table if exists rubrics;
create table rubrics
(
    id         varchar(42) PRIMARY KEY,
    title       varchar(255)                       null,
    description       text                      null
);
-- insert common rubric
insert into rubrics(id, title, description) VALUES ('00000000-0000-0000-00000000', 'Go for fun', 'Go rubric for Golang funs');

-- drop table if exists comments;
create table comments
(
    id         varchar(42) PRIMARY KEY,
    author_id       varchar(42)                       null,
    content       text                      not null,
    count_of_stars int default 0 not null,
    post_id varchar(42) not null
);

-- drop table if exists posts;
create table posts
(
    id         varchar(42) PRIMARY KEY,
    title       varchar(1000)                      not null,
    author_id       varchar(42)                       null,
    rubric_id varchar(42)                       null,
    tags    json null,
    state   SET('write', 'moderate', 'public', 'blocked'),
    content       text                      not null,
    created_at datetime default CURRENT_TIMESTAMP null,
    modified_at datetime default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    parent_post_id varchar(42)                       null,
    count_of_views int default 0 not null,
    count_of_stars int default 0 not null,
    comments_ids json null
);

alter table comments
add foreign key (post_id) references posts(id)
    on update cascade
    on delete cascade,
add foreign key (author_id) references users(id)
    on update cascade
    on delete set null;

alter table posts
add foreign key (rubric_id) references rubrics(id)
    on update cascade
    on delete set null,
add foreign key (author_id) references users(id)
    on update cascade
    on delete set null;