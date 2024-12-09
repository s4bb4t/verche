package config

import (
	"flag"
	"fmt"
	"os"
)

type fileSystem struct {
	BasePath          string
	PathToFile        string
	PathToVerchedFile string
}

type Config struct {
	GoVersion  string
	FileSystem fileSystem
	Method     byte
}

func newFS(path string) fileSystem {
	goModPath := path + "/go.mod"
	verchedGoModPath := path + "/verched_go.mod"

	return fileSystem{
		BasePath:          path,
		PathToFile:        goModPath,
		PathToVerchedFile: verchedGoModPath,
	}
}

const (
	manualUsage = "manual"
)

func MustLoad() *Config {
	inputPath := flag.String("in", "", "Path to the project directory containing go.mod (required)")
	goVersion := flag.String("v", "1.23.0", "Version of Golang")
	method := flag.String("m", manualUsage, "Method of usage")

	flag.Parse()

	if *inputPath == "" {
		fmt.Println("Error: Input path is required. Use the -in flag to specify the path.")
		flag.Usage()
		os.Exit(1)
	}

	if *goVersion == "" {
		fmt.Println("Error: Go version is required. Use the -v flag to specify the version.")
		flag.Usage()
		os.Exit(1)
	}

	var cfg Config
	cfg.GoVersion = *goVersion
	cfg.FileSystem = newFS(*inputPath)

	if *method == manualUsage {
		cfg.Method = 0
	} else {
		cfg.Method = 1
	}

	return &cfg
}
