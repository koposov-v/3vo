package main

import (
	"authjwt/internal/init/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Default().Fatal(err)
	}
}
