package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/smakimka/tack/internal/html"
	"github.com/smakimka/tack/internal/model"
	"github.com/smakimka/tack/internal/monitor/storage"
)

type EndpointHandler struct {
	s storage.Storage
}

func NewEndpointHandler(s storage.Storage) *EndpointHandler {
	return &EndpointHandler{s}
}

func (h *EndpointHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	html := html.New()

	tackName := r.PathValue("tackName")
	endpointType := r.PathValue("endpointType")
	endpointName := r.PathValue("endpointName")

	html.SetTitle(endpointName)

	var endpoint model.EndpointData
	switch endpointType {
	case "proxies":
		endpointData, ok := h.s.GetProxy(tackName, endpointName)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		endpoint = endpointData
	case "balancers":
		endpointData, ok := h.s.GetBalancer(tackName, endpointName)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		endpoint = endpointData
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body := []string{fmt.Sprintf("<h1>%s<h1>", endpointName)}
	body = append(body, fmt.Sprintf("<h2>TotalRecieved: %d bytes<h2>", endpoint.TotalRecieved))
	body = append(body, fmt.Sprintf("<h2>TotalSent: %d bytes<h2>", endpoint.TotalSent))

	body = append(body, "<h1>Clients:<h1>")
	for clientAddr, client := range endpoint.Clients {
		body = append(body, "<div>")
		body = append(body, fmt.Sprintf("<h4>%s</h4>", clientAddr))
		body = append(body, fmt.Sprintf("<h4>TotalConnections: %d<h4>", client.TotalConnections))
		body = append(body, fmt.Sprintf("<h4>TotalRecieved: %d bytes<h4>", client.TotalRecieved))
		body = append(body, fmt.Sprintf("<h4>TotalSent: %d bytes<h4>", client.TotalSent))

		body = append(body, "<h3>Servers data:<h3>")
		body = append(body, "<table border=1>")
		body = append(body, `
			<tr>
				<th>Server address</th>
				<th>Total connections</th>
				<th>Average connection duration</th>
				<th>Bytes recieved</th>
				<th>Bytes sent</th>
	 		</tr>`)

		for serverAddr, server := range client.Servers {
			body = append(body, fmt.Sprintf(`
				<tr>
					<th>%s</th>
					<th>%d</th>
					<th>%f</th>
					<th>%d</th>
					<th>%d</th>
				</tr>`, serverAddr, server.TotalConnections, server.AvgConnDuration, server.TotalRecieved, server.TotalSent))
		}

		body = append(body, "</table>")
		body = append(body, "</div>")
	}

	html.SetBody(strings.Join(body, "\n"))
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html.String()))
}
