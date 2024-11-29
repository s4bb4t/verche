package main

import (
	"github.com/s4bb4t/verche/pkg/config"
	"github.com/s4bb4t/verche/pkg/updater"
)

func main() {
	cfg := config.MustLoad()

	updater.Update("C:\\Users\\dmitriy.bratishkin\\GolandProjects\\auth", cfg.GoVersion)
}
