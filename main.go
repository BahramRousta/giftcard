package main

import (
	"fmt"
	"giftCard/internal/adaptor/giftcard"
	"log"
)

func main() {

	//server := api.NewServer()
	//
	//server.SetupRoutes()
	//
	//server.Start(":8000")

	gfc := adaptor.NewGiftCard()
	fmt.Println("secret", gfc.ClientSecret)
	token, err := gfc.Auth()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(token)
}
