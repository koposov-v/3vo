package main

import (
	"kitchen-queue/internal/init/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Default().Fatal(err)
	}
}
