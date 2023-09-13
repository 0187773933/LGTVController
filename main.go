package main

import (
	"fmt"
	"time"
	"net"
	"os"
	"io/ioutil"
	"encoding/json"
	websocket "github.com/gorilla/websocket"
)

func get_handshake_json() ( hand_shake_json interface{} ) {
	hand_shake_file , _ := os.Open( "handshake.json" )
	defer hand_shake_file.Close()
	hand_shake_bytes , _ := ioutil.ReadAll( hand_shake_file )
	json.Unmarshal( hand_shake_bytes , &hand_shake_json )
	return
}

type Payload map[string]string
type Message struct {
	Id string `json:"id"`
	Type string `json:"type"`
	Uri string `json:"uri"`
	Payload Payload `json:"payload"`
	ClientKey string `json:"client-key"`
}

// https://github.com/48723247842/LGTVController/tree/master/LGTVController
// https://github.com/TheRealLink/pylgtv/blob/master/pylgtv/webos_client.py#L92
func main() {

	tv_ip := "192.168.1.6"
	webscoket_port := "3000"
	websocket_url := fmt.Sprintf( "ws://%s:%s/" , tv_ip , webscoket_port )
	// hand_shake_json := get_handshake_json()
	dialer := websocket.DefaultDialer
	dialer.NetDial = func(network, addr string) (net.Conn, error) {
		conn, err := net.Dial(network, addr)
		if err != nil {
			return nil, err
		}
		// Set TCP_NODELAY option if required
		if tcpConn, ok := conn.(*net.TCPConn); ok {
			tcpConn.SetNoDelay(true)
		}
		return conn, nil
	}
	conn, _, err := dialer.Dial( websocket_url , nil )
	time.Sleep( 1 * time.Second )

	// Pairing if Never Paired Before
	// err = conn.WriteJSON( hand_shake_json )
	// if err != nil { panic( err ) }
	// _ , firstResponseBytes, err := conn.ReadMessage()
	// fmt.Println( string( firstResponseBytes ) )
	// // json.Unmarshal(firstResponseBytes, &first_response)
	// ioutil.WriteFile("pairing_one_response.json", firstResponseBytes, 0644)
	// _ , secondResponseBygtes , err := conn.ReadMessage()
	// fmt.Println( string( secondResponseBygtes ) )
	// ioutil.WriteFile( "pairing_two_response.json" , secondResponseBygtes , 0644 )

	// If Paired Before , then just Resend Handshake Data
	// err = conn.WriteJSON( hand_shake_json )
	// if err != nil { panic( err ) }
	// _ , firstResponseBytes, err := conn.ReadMessage()
	// fmt.Println( string( firstResponseBytes ) )

	// Example Command
	uri := fmt.Sprintf( "ssap://%s" , "audio/volumeUp" )
	request_type := "volume_up"
	command_count := "1"

	// "audio/volumeUp"
	payload := &Payload{}
	message := &Message{
		Id: fmt.Sprintf( "lgtv_%s_%s" , request_type , command_count ) ,
		Type: "request" ,
		Uri: uri ,
		Payload: *payload ,
		ClientKey: "asdf" ,
	}
	fmt.Println( message )
	err = conn.WriteJSON(  message )
	if err != nil { panic( err ) }
	_ , commandResponseBytes , _ := conn.ReadMessage()
	fmt.Println( string( commandResponseBytes ) )

}