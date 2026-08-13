package main

import (
	"log"
)

func main() {
	app, err := NewApp()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	err = app.Run()

	if err != nil {
		log.Fatalf("%+v", err)
	}
}
