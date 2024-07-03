package main

import (
	"log"
	"todo/internal/config"
	todoapp "todo/internal/todo-app"
)

func main() {
	config := config.NewEnvConfig()

	app := todoapp.NewApp(config)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
