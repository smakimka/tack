package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/smakimka/tack/internal/html"
	"github.com/smakimka/tack/internal/monitor/storage"
)

type TackHandler struct {
	s storage.Storage
}

func NewTackHandler(s storage.Storage) *TackHandler {
	return &TackHandler{s}
}

func (h *TackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	html := html.New()

	tackName := r.PathValue("tackName")
	html.SetTitle(tackName)

	tack, ok := h.s.GetTack(tackName)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body := []string{}

	if tack.Alive {
		body = append(body, "<h1>online<h1>")
	} else {
		body = append(body, "<h1>offline<h1>")
	}

	body = append(body, fmt.Sprintf("<h1>last update: %s<h1>", tack.LastUpdated.String()))

	body = append(body, "<h2>Proxies:</h2>\n<ul>")
	for proxyName := range tack.Proxies {
		body = append(body, fmt.Sprintf("<li><a href=\"/tacks/%s/proxies/%s\">%s</a></li>", tackName, proxyName, proxyName))
	}
	body = append(body, "</ul>")

	body = append(body, "<h2>Balancers:</h2>\n<ul>")
	for balancerName := range tack.Balancers {
		body = append(body, fmt.Sprintf("<li><a href=\"/tacks/%s/balancers/%s\">%s</a></li>", tackName, balancerName, balancerName))
	}
	body = append(body, "</ul>")

	html.SetBody(strings.Join(body, "\n"))

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html.String()))
}
