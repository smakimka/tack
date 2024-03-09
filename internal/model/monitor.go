package model

type EndpointData struct {
	TotalRecieved int                   `json:"recieved"`
	TotalSent     int                   `json:"sent"`
	Clients       map[string]ClientData `json:"clients"`
}

type ClientData struct {
	TotalRecieved    int                   `json:"recieved"`
	TotalSent        int                   `json:"sent"`
	TotalConnections int                   `json:"connections"`
	Servers          map[string]ServerData `json:"servers"`
}

type ServerData struct {
	TotalRecieved    int     `json:"recieved"`
	TotalSent        int     `json:"sent"`
	TotalConnections int     `json:"connections"`
	AvgConnDuration  float64 `json:"avg_conn_duration"`
}

type SendData struct {
	Name      string                  `json:"name"`
	Endpoints map[string]EndpointData `json:"endpoints"`
	Balancers map[string]EndpointData `json:"balancers"`
}
