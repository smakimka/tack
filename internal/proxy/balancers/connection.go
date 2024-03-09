package balancers

import (
	"context"
	"net"

	"github.com/rs/zerolog/log"

	"github.com/smakimka/tack/internal/model"
	"github.com/smakimka/tack/internal/proxy/monitor"
	"github.com/smakimka/tack/internal/proxy/tcp"
)

type ConnectionHandler interface {
	HandleConnection(ctx context.Context, clientConn net.Conn)
}

type SimpleConnectionHandler struct {
	monitorChan chan monitor.MonitorData
	splitter    Splitter
}

func NewSimpleConnectionHandler(monitorChan chan monitor.MonitorData, s Splitter) *SimpleConnectionHandler {
	return &SimpleConnectionHandler{monitorChan, s}
}

func (c *SimpleConnectionHandler) HandleConnection(ctx context.Context, clientConn net.Conn) {
	defer clientConn.Close()
	balancerName := ctx.Value(model.BalancerKey).(string)
	workerNum := ctx.Value(model.WorkerKey).(int)

	log.Debug().Int("worker_num", workerNum).Str("balancer_name", balancerName).Msg("serving new client connection")

	addr := c.splitter.Next()

	tcpServer, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Err(err).Int("worker_num", workerNum).Str("balancer_name", balancerName).Msg("error resolving tcp addr")
		return
	}

	destConn, err := net.Dial("tcp", tcpServer.String())
	if err != nil {
		log.Err(err).Int("worker_num", workerNum).Str("balancer_name", balancerName).Msg("error dialing destination")
		return
	}

	if c.monitorChan == nil {
		tcp.Proxy(ctx, clientConn, destConn)
	} else {
		tcp.ProxyWithMonitor(ctx, c.monitorChan, clientConn, destConn)
	}

	log.Info().Int("worker_num", workerNum).Str("balancer_name", balancerName).Msg("end of serving connection")
}
