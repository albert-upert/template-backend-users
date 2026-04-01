package main

import (
	"log"

	"github.com/albert-upert/template-backend-users/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
