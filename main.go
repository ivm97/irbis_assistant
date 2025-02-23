package main

import (
	"github.com/irbis_assistant/internal/app"
	"github.com/irbis_assistant/internal/config"
	"log/slog"
	"os"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}

	slgr := slog.New(slog.NewTextHandler(os.Stdout, nil))

	iass := app.New(slgr, cfg.BackendAddr, cfg.ClientPath)

	if err := iass.Start(); err != nil {
		panic(err)
	}

}
