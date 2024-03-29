package main

import (
	"log"

	"github.com/nemesisesq/vaux_server/actions"
)

func main() {
	app := actions.App()
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}
