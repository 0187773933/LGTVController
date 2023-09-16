package controller

import (
	"fmt"
	"time"
	"math/rand"
	"strconv"
	"net"
	"encoding/json"
	utils "github.com/0187773933/LGTVController/v1/utils"
	types "github.com/0187773933/LGTVController/v1/types"
	websocket "github.com/gorilla/websocket"
)

const letters = "abcdefghijklmnopqrstuvwxyz"

type Controller struct {
	Config *types.ConfigFile
	Endpoints types.Endpoints
	API_Data types.API
}

func New( config *types.ConfigFile ) ( ctrl *Controller ) {
	// endpoints := utils.GetEndpoints( "./v1/misc/endpoints.yaml" )
	endpoints := utils.GetEndpoints()
	api := make( types.API )
	ctrl = &Controller{
		Config: config ,
		Endpoints: endpoints ,
		API_Data: api ,
	}
	for endpoint_name := range endpoints {
		ctrl.API_Data[ endpoint_name ] = ctrl.SendCommand
	}
	return
}

func ( ctrl *Controller ) GenMessageID() ( result string ) {
	rand.Seed( time.Now().UnixNano() )
	prefix_a := make( []byte , 3 )
	for i := range prefix_a { prefix_a[ i ] = letters[ rand.Intn( len( letters ) ) ] }
	prefix := string( prefix_a )
	suffix_int := ( rand.Intn( 100 ) + 1 )
	suffix := strconv.Itoa( suffix_int )
	result = fmt.Sprintf( "lgtv_%s_%s" , prefix , suffix )
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

func ( ctrl *Controller ) SendCommand( endpoint types.Endpoint ) ( result string ) {
	fmt.Println( "SendCommand()" , endpoint )
	if ctrl.Config.ClientKey == "" {
		fmt.Println( "Client Key is Empty !!!" )
		fmt.Println( "You have to Pair , and Recieve a Client Key !!!" )
		return
	}
	if endpoint.Payload == nil { endpoint.Payload = types.Payload{} } // might not need this ?
	message_id := ctrl.GenMessageID()
	uri := fmt.Sprintf( "ssap://%s" , endpoint.Path )
	message := &types.Message{
		Id: message_id ,
		Type: "request" ,
		Uri: uri ,
		Payload: endpoint.Payload ,
		ClientKey: ctrl.Config.ClientKey ,
	}
	timeout := time.After( time.Duration( ctrl.Config.TimeoutSeconds ) * time.Second )
	message_channel := make( chan []byte )
	var ws *websocket.Conn
	go func() {
		ws = ctrl.Connect()
		// fmt.Println( "Sending :" , message )
		message_json , _ := json.Marshal( message )
		fmt.Println( "Sending :" , string( message_json ) )
		write_err := ws.WriteJSON( message )
		if write_err != nil { panic( write_err ) }
		_ , response_bytes , response_err := ws.ReadMessage()
		ws.Close()
		if response_err == nil {
			message_channel <- response_bytes
		} else {
			close( message_channel )
		}
	}()
	select {
		case response_bytes , ok := <-message_channel:
			if ok {
				result = string( response_bytes )
			} else {
				ws.Close()
				result = "error reading message"
			}
		case <-timeout:
			ws.Close()
			result = "timeout while reading message"
	}
	return
}

func ( ctrl *Controller ) API( endpoint_name string , payload ...types.Payload ) ( result string ) {
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
	result = ctrl.API_Data[ endpoint_name ]( new_endpoint )
	return
}

func ( ctrl *Controller ) WakeOnLan() {
	mac_bytes , _ := net.ParseMAC( ctrl.Config.TVMAC )
	magic_packet := []byte{}
	for i := 0; i < 6; i++ {
		magic_packet = append( magic_packet , 0xFF )
	}
	for i := 0; i < 16; i++ {
		magic_packet = append( magic_packet , mac_bytes... )
	}
	addr := &net.UDPAddr{
		IP:   net.IPv4bcast ,
		Port: 9 ,
	}
	conn , _ := net.DialUDP( "udp" , nil , addr )
	defer conn.Close()
	conn.Write( magic_packet )
}