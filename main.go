package main

import (
	"giftCard/api"
)

func main() {
	server := api.NewServer()
	server.SetupRoutes()
	server.Logger.Fatal(server.Start(":8000"))
}
