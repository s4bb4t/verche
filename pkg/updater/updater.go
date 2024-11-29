package updater

import (
	"bufio"
	"fmt"
	"golang.org/x/mod/semver"
	"os"
	"strings"

	"github.com/s4bb4t/verche/pkg/handler"
	"github.com/s4bb4t/verche/pkg/liner"
)

func Update(path string, goVersion string) {
	goModPath := path + "/go.mod"
	//newFilePath := path + "/verched_go.mod"
	newFilePath := "verched_go.mod"

	fmt.Printf("Processing file: %s\n\n", goModPath)

	file, err := os.Open(goModPath)
	if err != nil {
		panic(fmt.Sprintf("Error opening file: %v\n", err))
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(fmt.Sprintf("Error closing file: %v\n", err))
		}
	}()

	newFile, err := os.Create(newFilePath)
	if err != nil {
		panic(fmt.Sprintf("Error creating new file: %v\n", err))
		return
	}
	defer func() {
		if err := newFile.Close(); err != nil {
			panic(fmt.Sprintf("Error closing new file: %v\n", err))
		}
	}()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(newFile)

	for scanner.Scan() {
		line := scanner.Text()
		if pkg, ver, ok := liner.TakeALook(line); ok {
			resp, err := handler.ParseResponse(handler.SendPackageRequest(pkg))
			if err != nil {
				panic(fmt.Sprintf("Error fetching package info for %s: %v\n", pkg, err))
			}

			maxVer := ver
			for _, art := range resp.Artifacts {
				version := art.Go.Version
				if semver.IsValid(version) && semver.Compare(version, maxVer) > 0 && art.State.Status == "PERMITTED" {
					maxVer = version
				}
			}

			if maxVer != "" {
				newLine := fmt.Sprintf("%s %s", pkg, maxVer)
				fmt.Printf("%s %s --> Latest Version: %s\n", pkg, ver, maxVer)
				if _, err := writer.WriteString("\t" + newLine + "\n"); err != nil {
					panic(err)
				}
			} else {
				panic("PACKAGE IS NOT PERMITTED")
			}
		} else {
			if strings.Contains(line, "toolchain") {
				fmt.Println(line)
				line = "toolchain go1.22.0"
				fmt.Println(line)
			} else if strings.Contains(line, "go 1.") {
				fmt.Println(line)
				line = "go " + goVersion
				fmt.Println(line)
			}
			if _, err := writer.WriteString(line + "\n"); err != nil {
				panic(err)
			}
		}
	}

	if err := writer.Flush(); err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}

	fmt.Printf("\nVerched! Take a look at %s", newFilePath)
}
