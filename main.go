package main

import (
	"wormhole/pkg/auth"
	"wormhole/pkg/database"
	"wormhole/pkg/server"
)

func main() {
	auth.Init()

	database.Init()
	defer database.Close()

	server.Init()
}
