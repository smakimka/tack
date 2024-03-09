package proxy

import (
	"context"
	"net"

	"github.com/rs/zerolog/log"

	"github.com/smakimka/tack/internal/model"
)

func worker(ctx context.Context, workChan chan net.Conn, endChan chan string, connectionHandler ConnectionHandler) {
	endpointName := ctx.Value(model.EndpointKey).(string)
	workerNum := ctx.Value(model.WorkerKey).(int)

	log.Info().Int("worker_num", workerNum).Str("endpoint_name", endpointName).Msg("started worker")

	for {
		select {
		case <-ctx.Done():
			return
		case clientConn := <-workChan:
			connectionHandler.HandleConnection(ctx, clientConn)
			endChan <- clientConn.RemoteAddr().String()
		}
	}
}
