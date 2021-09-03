package main

import (
	"log"
	"xrate/app"
)

func main() {
	app := app.NewApi()
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
