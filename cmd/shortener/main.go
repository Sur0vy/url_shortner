package main

import (
	"github.com/Sur0vy/url_shortner.git/internal/server"
)

func main() {
	const port = 8080
	server.StartServer(port)
}
