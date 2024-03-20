package proxy

import (
	"bytes"
	"context"
	"io"
	"net"

	"github.com/rs/zerolog/log"

	"github.com/smakimka/tack/internal/model"
	"github.com/smakimka/tack/internal/proxy/http"
	"github.com/smakimka/tack/internal/proxy/monitor"
	"github.com/smakimka/tack/internal/proxy/tcp"
)

type ConnectionHandler interface {
	HandleConnection(ctx context.Context, clientConn net.Conn)
}

type SimpleConnectionHandler struct {
	monitorChan chan monitor.MonitorData
}

func NewSimpleConnectionHandler(monitorChan chan monitor.MonitorData) *SimpleConnectionHandler {
	return &SimpleConnectionHandler{monitorChan}
}

func (c *SimpleConnectionHandler) HandleConnection(ctx context.Context, clientConn net.Conn) {
	defer clientConn.Close()
	endpointName := ctx.Value(model.EndpointKey).(string)
	workerNum := ctx.Value(model.WorkerKey).(int)

	log.Debug().Int("worker_num", workerNum).Str("endpoint_name", endpointName).Msg("serving new client connection")

	request := http.NewRequest()
	if err := request.Parse(ctx, clientConn); err != nil {
		log.Err(err).Int("worker_num", workerNum).Str("endpoint_name", endpointName).Msg("error processing request")
		return
	}

	tcpServer, err := net.ResolveTCPAddr("tcp", request.Addr)
	if err != nil {
		log.Err(err).Int("worker_num", workerNum).Str("endpoint_name", endpointName).Msg("error resolving tcp addr")
		return
	}

	destConn, err := net.Dial("tcp", tcpServer.String())
	if err != nil {
		log.Err(err).Int("worker_num", workerNum).Str("endpoint_name", endpointName).Msg("error dialing destination")
		return
	}
	defer destConn.Close()

	io.Copy(destConn, bytes.NewReader(request.Buffer))

	if c.monitorChan == nil {
		tcp.Proxy(ctx, clientConn, destConn)
	} else {
		tcp.ProxyWithMonitor(ctx, c.monitorChan, clientConn, destConn)
	}

	log.Debug().Int("worker_num", workerNum).Str("endpoint_name", endpointName).Msg("end of serving connection")
}
