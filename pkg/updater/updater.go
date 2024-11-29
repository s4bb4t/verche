package updater

import (
	"bufio"
	"fmt"
	"golang.org/x/mod/semver"
	"io"
	"os"
	"strings"

	"github.com/s4bb4t/verche/pkg/handler"
	"github.com/s4bb4t/verche/pkg/liner"
)

func Update(path string, goVersion string) {
	goModPath := path + "/go.mod"
	newFilePath := path + "/verched_go.mod"

	err := copyFile(goModPath, path+"/go.mod-old")
	if err != nil {
		panic(fmt.Sprintf("Error copying file: %v\n", err))
	}

	fmt.Printf("Processing file: %s\n\n", goModPath)

	// Открываем исходный файл
	file, err := os.Open(goModPath)
	if err != nil {
		panic(fmt.Sprintf("Error opening file: %v\n", err))
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(fmt.Sprintf("Error closing file: %v\n", err))
		}
	}()

	// Создаем новый файл для записи
	newFile, err := os.Create(newFilePath)
	if err != nil {
		panic(fmt.Sprintf("Error creating new file: %v\n", err))
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
				line = "toolchain go1.22.0"
			} else if strings.Contains(line, "go 1.") {
				line = "go " + goVersion
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

	fmt.Println("\nOverwriting the original go.mod with the updated content.")
	overwriteFile(newFilePath, goModPath)

	fmt.Printf("\nVerched! Updated file: %s", goModPath)
}

func overwriteFile(sourceFile, destFile string) {
	src, err := os.Open(sourceFile)
	if err != nil {
		panic(fmt.Sprintf("Error opening source file: %v\n", err))
	}
	defer src.Close()

	dst, err := os.Create(destFile)
	if err != nil {
		panic(fmt.Sprintf("Error creating destination file: %v\n", err))
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		panic(fmt.Sprintf("Error copying content: %v\n", err))
	}
}

func copyFile(source, destination string) error {
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}
