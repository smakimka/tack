package http

import (
	"context"
	"errors"
	"net"
	"strings"

	"github.com/rs/zerolog/log"
)

var ErrParseRequest = errors.New("error parsing request")

type HttpRequest struct {
	Addr   string
	Buffer []byte
}

func NewRequest() *HttpRequest {
	return &HttpRequest{Buffer: make([]byte, 1024)}
}

func (r *HttpRequest) Parse(ctx context.Context, conn net.Conn) error {
	n, err := conn.Read(r.Buffer)
	if err != nil {
		return err
	}

	if n == 0 {
		return ErrParseRequest
	}
	r.Buffer = r.Buffer[:n]

	log.Info().Msg(string(r.Buffer))
	host, err := getHost(string(r.Buffer))
	if err != nil {
		return err
	}

	if !strings.Contains(host, ":") {
		r.Addr = host + ":80"
	} else {
		r.Addr = host
	}

	return nil
}

func getHost(buffer string) (string, error) {
	var host string

	split := strings.Split(buffer, "\n")
	if len(split) == 0 {
		return host, ErrParseRequest
	}

	for i := 1; i < len(split); i++ {
		lineSplit := strings.Split(split[i], ": ")
		if len(lineSplit) < 2 {
			return host, ErrParseRequest
		}

		if lineSplit[0] == "Host" {
			host = strings.TrimSpace(lineSplit[1])
			return host, nil
		}
	}

	return host, ErrParseRequest
}
