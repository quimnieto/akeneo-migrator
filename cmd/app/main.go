package main

import (
	"log"

	"akeneo-migrator/cmd/app/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
