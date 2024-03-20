package tcp

import (
	"context"
	"net"
	"time"

	"github.com/smakimka/tack/internal/model"
	"github.com/smakimka/tack/internal/proxy/monitor"
)

func chanFromConn(conn net.Conn) chan []byte {
	c := make(chan []byte)

	go func() {
		b := make([]byte, 1024)

		for {
			n, err := conn.Read(b)
			if n > 0 {
				res := make([]byte, n)
				copy(res, b[:n])
				c <- res
			}
			if err != nil {
				c <- nil
				break
			}
		}
	}()

	return c
}

func ProxyWithMonitor(ctx context.Context, monitorChan chan monitor.MonitorData, clientConn net.Conn, destConn net.Conn) {
	monitorData := monitor.MonitorData{
		ClientAddr: clientConn.RemoteAddr().String(),
		DestAddr:   destConn.RemoteAddr().String(),
	}

	endpointNameValue := ctx.Value(model.EndpointKey)
	if endpointNameValue != nil {
		endpointName := endpointNameValue.(string)

		monitorData.Type = string(model.EndpointKey)
		monitorData.Name = endpointName
	} else {
		balancerName := ctx.Value(model.BalancerKey).(string)

		monitorData.Type = string(model.BalancerKey)
		monitorData.Name = balancerName
	}

	start := time.Now()
	defer func() {
		monitorData.ConnDuration = time.Since(start)
		monitorChan <- monitorData
	}()

	clientChan := chanFromConn(clientConn)
	destChan := chanFromConn(destConn)

	for {
		select {
		case clientBytes := <-clientChan:
			if clientBytes == nil {
				return
			} else {
				monitorData.BytesRecievedFromClient += len(clientBytes)
				destConn.Write(clientBytes)
				monitorData.BytesSentToDest += len(clientBytes)
			}
		case destBytes := <-destChan:
			if destBytes == nil {
				return
			} else {
				monitorData.BytesRecievedFromDest += len(destBytes)
				clientConn.Write(destBytes)
				monitorData.BytesSentToClient += len(destBytes)
			}
		case <-ctx.Done():
			return
		}
	}
}

func Proxy(ctx context.Context, clientConn net.Conn, destConn net.Conn) {
	clientChan := chanFromConn(clientConn)
	destChan := chanFromConn(destConn)

	for {
		select {
		case clientBytes := <-clientChan:
			if clientBytes == nil {
				return
			} else {
				destConn.Write(clientBytes)
			}
		case destBytes := <-destChan:
			if destBytes == nil {
				return
			} else {
				clientConn.Write(destBytes)
			}
		case <-ctx.Done():
			return
		}
	}
}

func CheckServerConnection(conn net.Conn) bool {
	one := make([]byte, 1)
	_, err := conn.Write(one)
	return err == nil
}
