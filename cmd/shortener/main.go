package main

import (
	"github.com/Sur0vy/url_shortner.git/internal/server"
)

func main() {

	const port = 1230
	server.StartServer(port)
}
