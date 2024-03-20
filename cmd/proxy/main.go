package main

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/smakimka/tack/internal/model"
	"github.com/smakimka/tack/internal/proxy/balancers"
	"github.com/smakimka/tack/internal/proxy/config"
	"github.com/smakimka/tack/internal/proxy/monitor"
	"github.com/smakimka/tack/internal/proxy/proxy"
)

func main() {
	config, err := config.Read()
	if err != nil {
		log.Err(err).Msg("error reading config")
		os.Exit(1)
	}

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

	totalWorkers := 0
	for _, endpoint := range config.Endpoints {
		totalWorkers += endpoint.Workers
	}
	for _, balancer := range config.Balancers {
		totalWorkers += balancer.Workers
	}
	var monitorChan chan monitor.MonitorData
	if config.Monitor {
		monitorChan = make(chan monitor.MonitorData, totalWorkers)
	}

	ctx := context.Background()

	monitor := monitor.New(config.Name, config.MonitorAddr, monitorChan)
	go monitor.Monitor(ctx)

	errChan := make(chan error)
	for name, endpoint := range config.Endpoints {
		endpointCtx := context.WithValue(ctx, model.EndpointKey, name)
		go proxy.Serve(endpointCtx, errChan, monitorChan, &endpoint)
	}

	for name, balancer := range config.Balancers {
		totalWorkers += balancer.Workers
		endpointCtx := context.WithValue(ctx, model.EndpointKey, name)
		go balancers.Serve(endpointCtx, errChan, monitorChan, &balancer)
	}

	err = <-errChan
	log.Fatal().Err(err).Msg("error creating listener")
}
