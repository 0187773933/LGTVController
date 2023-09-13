package controller

import (
	"fmt"
	// "time"
	"net"
	types "github.com/0187773933/LGTVController/v1/types"
	websocket "github.com/gorilla/websocket"
)

type Controller struct {
	Config *types.ConfigFile
}

func New( config *types.ConfigFile ) ( ctrl *Controller ) {
	ctrl = &Controller{
		Config: config ,
	}
	return
}

func ( ctrl *Controller ) Connect() ( result *websocket.Conn ) {
	websocket_url := fmt.Sprintf( "ws://%s:%s/" , ctrl.Config.TVIP , ctrl.Config.WebSocketPort )
	dialer := websocket.DefaultDialer
	dialer.NetDial = func( network , addr string ) ( net.Conn , error ) {
		conn , err := net.Dial( network , addr )
		if err != nil { return nil , err }
		if tcp_conn , ok := conn.( *net.TCPConn ); ok {
			tcp_conn.SetNoDelay( true )
		}
		return conn , nil
	}
	conn , _ , dial_err := dialer.Dial( websocket_url , nil )
	if dial_err != nil { panic( dial_err ); }
	// defer conn.close()
	// time.Sleep( 1 * time.Second )
	result = conn
	return
}

func ( ctrl *Controller ) Pair() ( result string ) {
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
	return
}


func ( ctrl *Controller ) TestCommand() ( result string ) {
	ws := ctrl.Connect()
	uri := fmt.Sprintf( "ssap://%s" , "audio/volumeUp" )
	request_name := "volume_up"
	command_count := "1"
	payload := &types.Payload{}
	message := &types.Message{
		Id: fmt.Sprintf( "lgtv_%s_%s" , request_name , command_count ) ,
		Type: "request" ,
		Uri: uri ,
		Payload: *payload ,
		ClientKey: "asdf" ,
	}
	write_err := ws.WriteJSON(  message )
	if write_err != nil { panic( write_err ) }
	_ , response_bytes , _ := ws.ReadMessage()
	result = string( response_bytes )
	ws.Close()
	return
}



