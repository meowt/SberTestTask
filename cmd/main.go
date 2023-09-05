package main

import (
	"log"

	"github.com/meowt/SberTestTask/cmd/modules"
	"github.com/meowt/SberTestTask/internal/config"
	"github.com/meowt/SberTestTask/internal/router"
	"github.com/meowt/SberTestTask/internal/storage"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Panic(err)
	}

	s, err := storage.New(&cfg.Storage)
	if err != nil {
		log.Panic(err)
	}

	handlerModule := modules.Setup(s)

	r := router.Setup(handlerModule)

	if err := r.Run(cfg.Router.Address); err != nil {
		log.Panic(err)
	}
}
