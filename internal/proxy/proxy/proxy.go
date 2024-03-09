package proxy

import (
	"context"
	"net"

	"github.com/rs/zerolog/log"

	"github.com/smakimka/tack/internal/model"
	"github.com/smakimka/tack/internal/proxy/limiter"
	"github.com/smakimka/tack/internal/proxy/monitor"
)

func Serve(ctx context.Context, c chan error, monitorChan chan monitor.MonitorData, config *model.ProxyEndpoint) {
	endpointNameValue := ctx.Value(model.EndpointKey)
	endpointName := endpointNameValue.(string)

	l, err := net.Listen("tcp", config.Addr)
	if err != nil {
		c <- err
		return
	}
	defer l.Close()

	connectionHandler := NewSimpleConnectionHandler(monitorChan)
	rateLimiter := limiter.New(config.SpeedLimit, config.TotalLimit)

	endChan := make(chan string, config.Workers)
	go func() {
		for {
			addr := <-endChan
			rateLimiter.End(addr)
		}
	}()

	workChan := make(chan net.Conn, config.Workers)
	for i := range config.Workers {
		workerCtx := context.WithValue(ctx, model.WorkerKey, i)
		go worker(workerCtx, workChan, endChan, connectionHandler)
	}

	log.Info().Str("endpoint_name", endpointName).Str("addr", config.Addr).Msg("serving endpoint")

	for {
		clientConn, err := l.Accept()
		if err != nil {
			log.Err(err).Str("endpoint_name", endpointName).Msg("error accepting client connection")
			continue
		}
		if rateLimiter != nil {
			if !rateLimiter.Recieve(clientConn.RemoteAddr().String()) {
				clientConn.Close()
				continue
			}
		}

		workChan <- clientConn
	}
}
