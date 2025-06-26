package main

import (
	"log"
	"order-core/internal/init/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Default().Fatal(err)
	}
}
