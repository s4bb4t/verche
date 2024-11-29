package main

import (
	"fmt"
	"github.com/s4bb4t/verche/pkg/config"
	"github.com/s4bb4t/verche/pkg/updater"
)

func main() {
	cfg := config.MustLoad()

	if err := updater.Update(cfg); err != nil {
		panic(err)
	}

	fmt.Println("Update and tidy process completed successfully.")
}
