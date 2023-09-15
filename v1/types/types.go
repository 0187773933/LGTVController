package types

type ConfigFile struct {
	TVIP string `yaml:"tv_ip"`
	WebSocketPort string `yaml:"websocket_port"`
	ClientKey string `yaml:"client_key"`
	TimeoutSeconds int `yaml:"timeout_seconds"`
}

type Payload map[string]interface{}

type Endpoint struct {
	Path string `yaml:"path"`
	Payload Payload `yaml:"payload"`
}
type Endpoints map[string]Endpoint

type API map[string]func( Endpoint )( string )

type Message struct {
	Id string `json:"id"`
	Type string `json:"type"`
	Uri string `json:"uri"`
	Payload Payload `json:"payload"`
	ClientKey string `json:"client-key"`
}