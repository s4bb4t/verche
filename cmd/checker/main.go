package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/s4bb4t/verche/pkg/config"
	"github.com/s4bb4t/verche/pkg/updater"
)

func main() {
	inputPath := flag.String("in", "", "Path to the project directory containing go.mod")
	flag.Parse()

	if *inputPath == "" {
		fmt.Println("Error: Input path is required. Use the -in flag to specify the path.")
		flag.Usage()
		os.Exit(1)
	}

	cfg := config.MustLoad()

	updater.Update(*inputPath, cfg.GoVersion)
}
