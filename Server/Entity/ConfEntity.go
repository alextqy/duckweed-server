package entity

type ConfEntity struct {
	TcpPort          string `json:"tcp_port"`
	UdpPort          string `json:"udp_port"`
	Lang             string `json:"lang"`
	InitialSpaceSize string `json:"initial_space_size"`
}
