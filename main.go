package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/art-frela/blog/infra"
)

// @title Blog API
// @version 1.0
// @description This is a simple blog server.

// @contact.url https://github.com/art-frela
// @contact.email art.frela@gmail.com

// @BasePath /api/v1

func main() {
	connSTR := flag.String("c", os.Getenv("DATABASE_URL"), "MySQL/MongoDB connection string, format: user:password@tcp(host:port)/database OR mongodb://locaslhost:27017, or use ENV: export DATABASE_URL=root:master@tcp(localhost:3306)/blog?parseTime=true or export DATABASE_URL=mongodb://locaslhost:27017")
	countExamplePosts := flag.Int("n", 0, "Count of example posts for inserting to storage")
	clearStorage := flag.Bool("clear", false, "Clear storage when start app")
	httpPORT := flag.String("p", "localhost:8888", "Host:TCPPort for HTTP server")
	flag.Parse()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	server := infra.NewBlogServer(*connSTR, *countExamplePosts, *clearStorage)
	srv := server.Run(*httpPORT)

	// Waiting for SIGINT (pkill -2)
	<-stop
	server.Stop(srv)
}
