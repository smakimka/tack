package router

import (
	"net/http"

	"github.com/smakimka/tack/internal/monitor/handlers"
	"github.com/smakimka/tack/internal/monitor/storage"
)

func New(s storage.Storage) *http.ServeMux {
	rootHandler := handlers.NewRootHandler(s)
	updateHandler := handlers.NewUpdateHandler(s)
	tackHandler := handlers.NewTackHandler(s)
	endpointHandler := handlers.NewEndpointHandler(s)

	mux := http.NewServeMux()

	mux.Handle("GET /", rootHandler)
	mux.Handle("GET /tacks/{tackName}/", tackHandler)
	mux.Handle("GET /tacks/{tackName}/{endpointType}/{endpointName}/", endpointHandler)

	mux.Handle("POST /api/update/", updateHandler)

	return mux
}
