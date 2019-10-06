package main

import (
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

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	server := infra.NewBlogServer(0, false)
	server.Run()

	// Waiting for SIGINT (pkill -2)
	<-stop
	server.Stop()
}
