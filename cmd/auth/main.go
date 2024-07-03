package main

import (
	"log"
	"todo/internal/auth"
	"todo/internal/config"
)

func main() {
	config := config.NewEnvConfig()

	server := auth.NewAuthServer(*config)

	err := server.Start()
	if err != nil {
		log.Fatal(err)
	}

}
