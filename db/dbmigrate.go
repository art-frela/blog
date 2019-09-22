package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const (
	createDB = `
	
	create table users
	(
		id         varchar(42) PRIMARY KEY,
		username       varchar(255)                       null,
		nick       varchar(255)                       null,
		email   varchar(500)                               null,
		created_at datetime default CURRENT_TIMESTAMP null,
		modified_at datetime default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
		user_role int default -1 not null,
		salt varchar(25) default 'saltsalt' not null
	);
	
	-- drop table if exists rubrics;
	create table rubrics
	(
		id         varchar(42) PRIMARY KEY,
		title       varchar(255)                       null,
		description       text                      null
	);
	
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
		on delete set null;`
)

func main() {
	//linux shell$ export DATABASE_URL=root:master@tcp(localhost:3306)/blog
	query, err := ioutil.ReadFile("blog.sql")
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Starting database 'blog' creation")
	for _, dbname := range []string{"DATABASE_URL"} {
		db, err := sql.Open("mysql", os.Getenv(dbname))
		if err != nil {
			log.Fatal(err)
		}
		if _, err = db.Query(string(query)); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Finished database 'blog' creation")
}
