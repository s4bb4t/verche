package main

import (
	"flag"
	"fmt"
	"github.com/s4bb4t/verche/pkg/config"
	"github.com/s4bb4t/verche/pkg/updater"
	"os"
	"os/exec"
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

	err := updater.CopyFile(*inputPath, *inputPath+"go.mod-old")
	if err != nil {
		panic(fmt.Sprintf("Error copying file: %v\n", err))
	}

	for i := 0; i < 2; i++ {
		fmt.Printf("Update and tidy iteration %d\n", i+1)

		updater.Update(*inputPath, cfg.GoVersion)

		err := runGoModTidy(*inputPath)
		if err != nil {
			fmt.Printf("Error running 'go mod tidy': %v\n", err)
			os.Exit(1)
		}
	}

	if err := os.Remove(*inputPath + "/verched_go.mod"); err != nil {
		panic(fmt.Sprintf("Error removing file: %v\n", err))
	}

	fmt.Println("Update and tidy process completed successfully.")
}

func runGoModTidy(path string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Running 'go mod tidy' in directory: %s\n", path)
	return cmd.Run()
}
