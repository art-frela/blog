package main

import (
	"github.com/art-frela/blog/infra"
)

func main() {
	server := infra.NewBlogServer("root:master@tcp(localhost:3306)/blog")
	server.Run(":8888")
}
