package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	queries = [][]string{
		{"drop", `drop database if exists blog;`},
		{"createDB", `create database blog;`},
		{"users", `create table blog.users
					(
						id         varchar(42)     PRIMARY KEY,
						username   varchar(255)    null,
						nick       varchar(255)    null,
						email      varchar(500)    null,
						created_at datetime        default CURRENT_TIMESTAMP null,
						modified_at datetime        default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
						user_role  int             default -1 not null,
						salt       varchar(25)     default 'saltsalt' not null,
						avatar_url varchar(512)    null
					);`},
		{"insertDefaultUser", `insert into blog.users (id, username, nick, email, avatar_url) VALUES ('00000000-0000-0000-00000000', 'anonimous', 'anonimous', 'user@example.com', 'https://getuikit.com/docs/images/avatar.jpg');`},
		{"rubrics", `create table blog.rubrics
					(
						id         varchar(42) PRIMARY KEY,
						title       varchar(255)                       null,
						description       text                      null
					);`},
		{"insertDefaultRubric", `insert into blog.rubrics(id, title, description) VALUES ('00000000-0000-0000-00000000', 'Go for fun', 'Go rubric for Golang funs');`},
		{"comments", `create table blog.comments
					(
						id         varchar(42) PRIMARY KEY,
						author_id       varchar(42)                       null,
						content       text                      not null,
						count_of_stars int default 0 not null,
						post_id varchar(42) not null
					);`},
		{"posts", `create table blog.posts
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
				);`},
		{"foreignKeycomments", `alter table blog.comments
								add foreign key (post_id) references blog.posts(id)
									on update cascade
									on delete cascade,
								add foreign key (author_id) references blog.users(id)
									on update cascade
									on delete set null;`},
		{"foreignKeyPosts", `alter table blog.posts
							add foreign key (rubric_id) references blog.rubrics(id)
								on update cascade
								on delete set null,
							add foreign key (author_id) references blog.users(id)
								on update cascade
								on delete set null;`},
	}
)

func main() {
	//linux shell$ export DATABASE_URL=root:master@tcp(localhost:3306)/blog
	// query, err := ioutil.ReadFile("blog.sql")
	// if err != nil {
	// 	log.Panic(err)
	// }

	fmt.Println("Starting database 'blog' creation")
	for _, dbname := range []string{"DATABASE_URL"} {
		db, err := sql.Open("mysql", os.Getenv(dbname))
		if err != nil {
			log.Fatal(err)
		}
		for _, queryPair := range queries {
			fmt.Println(queryPair[0])
			if _, err = db.Exec(queryPair[1]); err != nil {
				log.Fatal(err)
			}
		}

	}
	fmt.Println("Finished database 'blog' creation")
}
