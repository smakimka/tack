package main

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/smakimka/tack/internal/monitor/config"
	"github.com/smakimka/tack/internal/monitor/router"
	"github.com/smakimka/tack/internal/monitor/storage"
)

func main() {
	config := config.Read()
	storage := storage.NewMemStorage()
	router := router.New(storage)

	switch config.Logginglevel {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "err":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		log.WithLevel(zerolog.ErrorLevel).Msg("error setting logging level")
	}

	aliveTicker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			<-aliveTicker.C
			go storage.CheckAlive()
		}
	}()

	if err := http.ListenAndServe(config.Addr, router); err != nil {
		log.Err(err).Msg("error starting")
	}
}
