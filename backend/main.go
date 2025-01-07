package main

import (
	"log"
	"workly/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal("Application failed to start: ", err)
	}
}
