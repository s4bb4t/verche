package updater

import (
	"bufio"
	"fmt"
	"os"

	"github.com/s4bb4t/verche/pkg/handler"
	"github.com/s4bb4t/verche/pkg/liner"
)

// Update processes a go.mod file to find and update package versions.
func Update(path string) {
	goModPath := path + "/go.mod"
	fmt.Printf("Processing file: %s\n", goModPath)

	file, err := os.Open(goModPath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	newFile, err := os.Create("verched_go.mod")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer newFile.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if pkg, ver, ok := liner.TakeALook(line); ok {
			resp, err := handler.ParseResponse(handler.SendPackageRequest(pkg))
			if err != nil {
				fmt.Printf("Error fetching package info for %s: %v\n", pkg, err)
				continue
			}

			maxVer := "v0.0.0"
			for _, art := range resp.Artifacts {
				if art.Go.Version > maxVer && art.State.Status == "PERMITTED" {
					maxVer = art.Go.Version
				}
			}

			fmt.Printf("%s, Latest Version: %s <- current version %s\n", pkg, maxVer, ver)
			_, err = newFile.WriteString(pkg + " " + maxVer + "\n")
			if err != nil {
				fmt.Printf("Error updating for %s: %v\n", pkg, err)
				return
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
}
