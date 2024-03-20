package balancers

import (
	"context"
	"net"

	"github.com/rs/zerolog/log"

	"github.com/smakimka/tack/internal/model"
	"github.com/smakimka/tack/internal/proxy/limiter"
	"github.com/smakimka/tack/internal/proxy/monitor"
)

func Serve(ctx context.Context, c chan error, monitorChan chan monitor.MonitorData, config *model.Balancer) {
	balancerNameValue := ctx.Value(model.EndpointKey)
	balancerName := balancerNameValue.(string)

	l, err := net.Listen("tcp", config.Addr)
	if err != nil {
		c <- err
		return
	}
	defer l.Close()

	var splitter Splitter
	if config.Type == "round_robin" {
		splitter = NewRoundRobinSplitter()
		err := splitter.AddAddrs(config.Addrs)
		if err != nil {
			c <- err
			return
		}
	}

	connectionHandler := NewSimpleConnectionHandler(monitorChan, splitter)
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

	log.Info().Str("balancer_name", balancerName).Msg("serving balancer")

	for {
		clientConn, err := l.Accept()
		if err != nil {
			log.Err(err).Str("balancer_name", balancerName).Msg("error accepting client connection")
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
