package main

import (
	"wormhole/pkg/config"
	"wormhole/pkg/database"
	"wormhole/pkg/email"
	"wormhole/pkg/jwt"
	"wormhole/pkg/server"
)

func main() {
	config.Init()
	jwt.Init()
	email.Init()
	database.Init()
	defer database.Close()
	server.Init()
}
