package main

import (
	handlers "github.com/AldoRuizP/proxy-app/api/handlers"
	server "github.com/AldoRuizP/proxy-app/api/server"
	utils "github.com/AldoRuizP/proxy-app/api/utils"
)

func main() {
	utils.LoadEnv()
	app := server.SetUp()
	handlers.HandleRedirection(app)
	server.RunServer(app)
}
