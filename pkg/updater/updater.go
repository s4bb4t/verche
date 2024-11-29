package updater

import (
	"bufio"
	"fmt"
	"github.com/s4bb4t/verche/pkg/config"
	"golang.org/x/mod/semver"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/s4bb4t/verche/pkg/handler"
	"github.com/s4bb4t/verche/pkg/liner"
)

func Update(cfg *config.Config) error {
	if err := update(cfg); err != nil {
		return err
	}
	if err := update(cfg); err != nil {
		return err
	}
	if err := os.Remove(cfg.FileSystem.PathToVerchedFile); err != nil {
		return err
	}
	return nil
}

func update(cfg *config.Config) error {
	file, err := os.Open(cfg.FileSystem.PathToFile)
	if err != nil {
		return err
	}
	defer file.Close()

	newFile, err := os.Create(cfg.FileSystem.PathToVerchedFile)
	if err != nil {
		return err
	}
	defer newFile.Close()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(newFile)

	for scanner.Scan() {
		line := scanner.Text()
		if pkg, ver, ok := liner.TakeALook(line); ok {
			resp, err := handler.ParseResponse(handler.SendPackageRequest(pkg))
			if err != nil {
				return fmt.Errorf("Error fetching package info for %s: %v", pkg, err)
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
					return err
				}
			} else {
				return fmt.Errorf("PACKAGE IS NOT PERMITTED: %s", pkg)
			}
		} else {
			if strings.Contains(line, "toolchain") {
				line = "toolchain go1.22.0"
			} else if strings.Contains(line, "go 1.") {
				line = "go " + cfg.GoVersion
			}
			if _, err := writer.WriteString(line + "\n"); err != nil {
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	if err := writer.Flush(); err != nil {
		return err
	}

	overwriteFile(cfg.FileSystem.PathToVerchedFile, cfg.FileSystem.PathToFile)
	return runGoModTidy(cfg.FileSystem.BasePath)
}

func overwriteFile(sourceFile, destFile string) {
	src, err := os.Open(sourceFile)
	if err != nil {
		panic(err)
	}
	defer src.Close()

	dst, err := os.Create(destFile)
	if err != nil {
		panic(err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		panic(err)
	}
}

func runGoModTidy(path string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Running 'go mod tidy' in directory: %s\n", path)
	return cmd.Run()
}
