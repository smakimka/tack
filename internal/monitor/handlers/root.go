package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/smakimka/tack/internal/html"
	"github.com/smakimka/tack/internal/monitor/storage"
)

type RootHandler struct {
	s storage.Storage
}

func NewRootHandler(s storage.Storage) *RootHandler {
	return &RootHandler{s}
}

func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	html := html.New()
	html.SetTitle("Tacks")

	body := []string{}

	for _, tack := range h.s.GetTacks() {
		body = append(body, fmt.Sprintf("<h1><a href=\"tacks/%s\">%s</a></h1>", tack, tack))
	}

	html.SetBody(strings.Join(body, "\n"))

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html.String()))
}
