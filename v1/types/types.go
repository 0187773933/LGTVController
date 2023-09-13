package types

type ConfigFile struct {
	TVIP string `yaml:"tv_ip"`
	WebSocketPort string `yaml:"websocket_port"`
	ClientKey string `yaml:"client_key"`
}

type Payload map[string]interface{}

type Endpoint struct {
	Path string `yaml:"path"`
	Payload Payload `yaml:"payload"`
}
type Endpoints map[string]Endpoint

type API map[string]func( Endpoint )

type Message struct {
	ClientKey string `json:"client-key"`
	Id string `json:"id"`
	Type string `json:"type"`
	Uri string `json:"uri"`
	Payload Payload `json:"payload"`
}