package monitor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/smakimka/tack/internal/model"
)

type MonitorData struct {
	Name                    string
	Type                    string
	ClientAddr              string
	DestAddr                string
	BytesRecievedFromClient int
	BytesRecievedFromDest   int
	BytesSentToDest         int
	BytesSentToClient       int
	ConnDuration            time.Duration
}

type Monitor struct {
	m             sync.Mutex
	c             chan MonitorData
	uuid          string
	client        http.Client
	addr          string
	name          string
	totalRecieved int
	totalSent     int
	endpointData  map[string]model.EndpointData
	balancerData  map[string]model.EndpointData
}

func New(name string, addr string, monitorChan chan MonitorData) *Monitor {
	return &Monitor{
		m:            sync.Mutex{},
		c:            monitorChan,
		client:       http.Client{},
		name:         name,
		addr:         addr,
		endpointData: map[string]model.EndpointData{},
		balancerData: map[string]model.EndpointData{},
	}
}

func (m *Monitor) Monitor(ctx context.Context) {
	go m.collectData(ctx)

	sendTicker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-sendTicker.C:
			go m.sendData()
		}
	}
}

func (m *Monitor) collectData(ctx context.Context) {
	log.Info().Msg("started collecting data")
	for {
		select {
		case <-ctx.Done():
			return
		case monitorData := <-m.c:
			log.Trace().Msg("got data from worker")
			m.writeData(&monitorData)
		}
	}
}

func (m *Monitor) writeData(data *MonitorData) {
	m.m.Lock()
	defer m.m.Unlock()

	m.totalRecieved += data.BytesRecievedFromClient
	m.totalRecieved += data.BytesRecievedFromDest
	m.totalSent += data.BytesSentToClient
	m.totalSent += data.BytesSentToDest

	if data.Type == string(model.EndpointKey) {
		endpointData, ok := m.endpointData[data.Name]
		if !ok {
			endpointData = model.EndpointData{Clients: map[string]model.ClientData{}}
		}
		m.endpointData[data.Name] = addEndpointData(endpointData, data)
	} else {
		endpointData, ok := m.balancerData[data.Name]
		if !ok {
			endpointData = model.EndpointData{Clients: map[string]model.ClientData{}}
		}
		m.balancerData[data.Name] = addEndpointData(endpointData, data)
	}
}

func addEndpointData(endpointData model.EndpointData, data *MonitorData) model.EndpointData {
	endpointData.TotalRecieved += data.BytesRecievedFromClient
	endpointData.TotalRecieved += data.BytesRecievedFromDest
	endpointData.TotalSent += data.BytesSentToClient
	endpointData.TotalSent += data.BytesSentToDest

	clientData, ok := endpointData.Clients[data.ClientAddr]
	if !ok {
		clientData = model.ClientData{Servers: map[string]model.ServerData{}}
	}

	clientData.TotalConnections++
	clientData.TotalRecieved += data.BytesRecievedFromClient
	clientData.TotalRecieved += data.BytesRecievedFromDest
	clientData.TotalSent += data.BytesSentToClient
	clientData.TotalSent += data.BytesSentToDest

	serverData, ok := clientData.Servers[data.DestAddr]
	if !ok {
		serverData = model.ServerData{}
	}

	serverData.TotalConnections++
	serverData.TotalRecieved += data.BytesRecievedFromClient
	serverData.TotalRecieved += data.BytesRecievedFromDest
	serverData.TotalSent += data.BytesSentToClient
	serverData.TotalSent += data.BytesSentToDest
	serverData.AvgConnDuration = (serverData.AvgConnDuration*float64(serverData.TotalConnections-1) + data.ConnDuration.Seconds()) / float64(serverData.TotalConnections)

	clientData.Servers[data.DestAddr] = serverData
	endpointData.Clients[data.ClientAddr] = clientData

	return endpointData
}

func (m *Monitor) sendData() {
	log.Trace().Msg("sending monitor data")

	m.m.Lock()
	data, err := json.Marshal(model.SendData{
		Name:      m.name,
		Endpoints: m.endpointData,
		Balancers: m.balancerData,
	})
	m.m.Unlock()

	if err != nil {
		log.Err(err).Msg("error marshaling send data")
		return
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/api/update/", m.addr), bytes.NewReader(data))
	if err != nil {
		log.Err(err).Msg("error creating request")
		return
	}

	resp, err := m.client.Do(req)
	if err != nil {
		log.Err(err).Msg("error sending request")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn().Int("status", resp.StatusCode).Msg("got not ok status sending monitor data")
		return
	}
	log.Trace().Msg("sending data ok")
}
