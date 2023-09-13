package types

type ConfigFile struct {
	TVIP string `yaml:"tv_ip"`
	WebSocketPort string `yaml:"websocket_port"`
	ClientKey string `yaml:"client_key"`
}

type Payload map[string]string

type Message struct {
	Id string `json:"id"`
	Type string `json:"type"`
	Uri string `json:"uri"`
	Payload Payload `json:"payload"`
	ClientKey string `json:"client-key"`
}