package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
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
