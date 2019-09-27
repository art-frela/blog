package main

import (
	"flag"
	"os"

	"github.com/art-frela/blog/infra"
)

func main() {
	connSTR := flag.String("c", os.Getenv("DATABASE_URL"), "MySQL connection string, format: user:password@tcp(host:port)/database, or use ENV: export DATABASE_URL=root:master@tcp(localhost:3306)/blog?parseTime=true")
	countExamplePosts := flag.Int("n", 0, "Count of example posts for inserting to storage")
	clearStorage := flag.Bool("clear", false, "Clear storage when start app")
	httpPORT := flag.String("p", ":8888", "Host:TCPPort for HTTP server")
	flag.Parse()

	server := infra.NewBlogServer(*connSTR, *countExamplePosts, *clearStorage)
	server.Run(*httpPORT)
}
