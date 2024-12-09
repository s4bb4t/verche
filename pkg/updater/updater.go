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
	"time"

	"github.com/s4bb4t/verche/pkg/handler"
	"github.com/s4bb4t/verche/pkg/liner"
)

const (
	manual byte = iota
	auto
)

func Update(cfg *config.Config) error {
	if cfg.Method == manual {
		if err := updateManual(cfg); err != nil {
			return err
		}
	} else if cfg.Method == auto {
		if err := updateAuto(cfg); err != nil {
			return err
		}
		if err := updateAuto(cfg); err != nil {
			return err
		}
		if err := os.Remove(cfg.FileSystem.PathToVerchedFile); err != nil {
			return err
		}
	}
	return nil
}

func updateAuto(cfg *config.Config) error {
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
				return fmt.Errorf("error fetching package info for %s: %v", pkg, err)
			}

			var lastDate time.Time
			var lastRequestedArtifact handler.Artifact
			currentVer := "v0.0.0"

			for _, art := range resp.Artifacts {
				// skip for rejected other packages that contain current package's name
				if art.State.Status != "PERMITTED" || art.Go.Name != pkg {
					continue
				}

				version := art.Go.Version

				if semver.IsValid(version) && semver.Compare(version, currentVer) > 0 {
					currentVer = version
					continue
				}

				date, err := time.Parse("02-01-2006 15:04:05.000 MST", art.State.RequestTime)
				if err != nil {
					return fmt.Errorf("error in request time: %w", err)
				}
				if date.After(lastDate) {
					lastDate = date
					lastRequestedArtifact = art
				}
			}
			if currentVer == "v0.0.0" {
				if lastRequestedArtifact.Go.Version == "" {
					return fmt.Errorf("\n\npackage not found\n\n")
				}
				currentVer = lastRequestedArtifact.Go.Version
			}

			if ver != currentVer {
				fmt.Println("changed", pkg, ver, currentVer)
			}
			newLine := fmt.Sprintf("%s %s", pkg, currentVer)
			if _, err := writer.WriteString("\t" + newLine + "\n"); err != nil {
				return err
			}
		} else {
			if strings.Contains(line, "go 1.") {
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

func updateManual(cfg *config.Config) error {
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
		if pkg, ver, ok := liner.TakeALook(line); ok && !strings.Contains(line, "// indirect") {
			resp, err := handler.ParseResponse(handler.SendPackageRequest(pkg))
			if err != nil {
				return fmt.Errorf("error fetching package info for %s: %v", pkg, err)
			}

			var lastDate time.Time
			var lastRequestedArtifact handler.Artifact
			currentVer := "v0.0.0"

			for _, art := range resp.Artifacts {
				// skip for rejected other packages that contain current package's name
				if art.State.Status != "PERMITTED" || art.Go.Name != pkg {
					continue
				}

				version := art.Go.Version
				if semver.IsValid(version) && semver.Compare(version, currentVer) > 0 {
					currentVer = version
					continue
				}

				date, err := time.Parse("02-01-2006 15:04:05.000 MST", art.State.RequestTime)
				if err != nil {
					return fmt.Errorf("error in request time: %w", err)
				}
				if date.After(lastDate) {
					lastDate = date
					lastRequestedArtifact = art
				}
			}
			if currentVer == "v0.0.0" {
				if lastRequestedArtifact.Go.Version == "" {
					return fmt.Errorf("\n\npackage not found\n\n")
				}
				currentVer = lastRequestedArtifact.Go.Version
			}

			newLine := fmt.Sprintf("%s %s", pkg, currentVer)
			fmt.Printf("%s %s --> Latest Version: %s\n", pkg, ver, currentVer)
			if _, err := writer.WriteString("\t" + newLine + "\n"); err != nil {
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

	return nil
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

	fmt.Println("go mod tidy")
	return cmd.Run()
}
