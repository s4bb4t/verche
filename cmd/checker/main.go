package main

import (
	"fmt"
	"github.com/s4bb4t/verche/pkg/config"
	"github.com/s4bb4t/verche/pkg/updater"
	"log"
)

func main() {
	cfg := config.MustLoad()

	if err := updater.Update(cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Update completed successfully!")
}
