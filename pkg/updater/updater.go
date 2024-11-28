package updater

import (
	"bufio"
	"fmt"
	"os"

	"github.com/s4bb4t/verche/pkg/handler"
	"github.com/s4bb4t/verche/pkg/liner"
)

func Update(path string) {
	goModPath := path + "/go.mod"
	newFilePath := path + "/verched_go.mod"
	fmt.Printf("Processing file: %s\n", goModPath)

	file, err := os.Open(goModPath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}()

	newFile, err := os.Create(newFilePath)
	if err != nil {
		fmt.Printf("Error creating new file: %v\n", err)
		return
	}
	defer func() {
		if err := newFile.Close(); err != nil {
			fmt.Printf("Error closing new file: %v\n", err)
		}
	}()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(newFile)

	for scanner.Scan() {
		line := scanner.Text()
		if pkg, ver, ok := liner.TakeALook(line); ok {
			resp, err := handler.ParseResponse(handler.SendPackageRequest(pkg))
			if err != nil {
				fmt.Printf("Error fetching package info for %s: %v\n", pkg, err)
				writer.WriteString(line + "\n") // Write original line in case of error
				continue
			}

			maxVer := "v0.0.0"
			for _, art := range resp.Artifacts {
				if art.Go.Version > maxVer && art.State.Status == "PERMITTED" {
					maxVer = art.Go.Version
				}
			}

			if maxVer != "v0.0.0" {
				newLine := fmt.Sprintf("%s %s", pkg, maxVer)
				fmt.Printf("%s, Latest Version: %s <- current version %s\n", pkg, maxVer, ver)
				writer.WriteString(newLine + "\n")
			} else {
				writer.WriteString(line + "\n")
			}
		} else {
			writer.WriteString(line + "\n")
		}
	}

	if err := writer.Flush(); err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
}
