package controller

import (
	"fmt"
	// "time"
	"net"
	utils "github.com/0187773933/LGTVController/v1/utils"
	types "github.com/0187773933/LGTVController/v1/types"
	websocket "github.com/gorilla/websocket"
)

type Controller struct {
	Config *types.ConfigFile
	Endpoints types.Endpoints
	API_Data types.API
}

func New( config *types.ConfigFile ) ( ctrl *Controller ) {
	endpoints := utils.GetEndpoints( "./v1/misc/endpoints.yaml" )
	api := make( types.API )
	ctrl = &Controller{
		Config: config ,
		Endpoints: endpoints ,
		API_Data: api ,
	}
	// for endpoint_name , endpoint := range endpoints {
	// 	x_endpoint := endpoint
	// 	ctrl.API[ endpoint_name ] = func() {
	// 		ctrl.SendCommand( x_endpoint )
	// 	}
	// }
	for endpoint_name := range endpoints {
		ctrl.API_Data[ endpoint_name ] = ctrl.SendCommand
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
	ws := ctrl.Connect()
	hand_shake_json := utils.GetHandshakeData()
	part_one_write_err := ws.WriteJSON( hand_shake_json )
	if part_one_write_err != nil { panic( part_one_write_err ) }
	_ , part_one_response_bytes , part_one_read_err := ws.ReadMessage()
	if part_one_read_err != nil { panic( part_one_read_err ) }
	result = string( part_one_response_bytes )
	fmt.Println( result )
	fmt.Println( "Use Your Remote to Accept Permission PopUp" )
	_ , part_two_response_bytes , part_two_read_err := ws.ReadMessage()
	if part_two_read_err != nil { panic( part_two_read_err ) }
	result = string( part_two_response_bytes )
	fmt.Println( result )
	fmt.Println( "save 'client-id' in config.yaml" )
	ws.Close()
	return
}

func ( ctrl *Controller ) SendCommand( endpoint types.Endpoint ) {
	fmt.Println( "SendCommand()" , endpoint )
	// if ctrl.Config.ClientKey == "" {
	// 	fmt.Println( "Client Key is Empty !!!" )
	// 	fmt.Println( "You have to Pair , and Recieve a Client Key !!!" )
	// 	return
	// }
	// ws := ctrl.Connect()

	// ws.Close()
	return
}


func ( ctrl *Controller ) API( endpoint_name string , payload ...types.Payload ) {
	original_endpoint := ( ctrl.Endpoints )[ endpoint_name ]
	new_payload := make( types.Payload )
	for k , v := range original_endpoint.Payload {
		new_payload[ k ] = v
	}
	if len( payload ) > 0 {
		for k , v := range payload[ 0 ] {
			new_payload[ k ] = v
		}
	}
	new_endpoint := types.Endpoint{
		Path: original_endpoint.Path ,
		Payload: new_payload ,
	}
	ctrl.API_Data[ endpoint_name ]( new_endpoint )
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



