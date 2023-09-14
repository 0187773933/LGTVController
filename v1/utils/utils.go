package utils

import (
	_ "embed"
	"io/ioutil"
	"encoding/json"
	yaml "gopkg.in/yaml.v2"
	types "github.com/0187773933/LGTVController/v1/types"
)

//go:embed misc/endpoints.yaml
var endpoints_yaml []byte

//go:embed misc/handshake.json
var handshake_json []byte

func GetConfig( file_path string ) ( result types.ConfigFile ) {
	file , file_read_err := ioutil.ReadFile( file_path )
	if file_read_err != nil { panic( file_read_err ) }
	error := yaml.Unmarshal( file , &result )
	if error != nil { panic( error ) }
	return
}

// func GetEndpoints( file_path string ) ( result types.Endpoints ) {
// 	file , _ := ioutil.ReadFile( file_path )
// 	error := yaml.Unmarshal( file , &result )
// 	if error != nil { panic( error ) }
// 	return
// }

func GetEndpoints() ( result types.Endpoints ) {
	err := yaml.Unmarshal( endpoints_yaml , &result )
	if err != nil { panic( err ) }
	return
}

// func GetHandshakeData() ( hand_shake_json interface{} ) {
// 	hand_shake_file , _ := os.Open( "./v1/utils/handshake.json" )
// 	defer hand_shake_file.Close()
// 	hand_shake_bytes , _ := ioutil.ReadAll( hand_shake_file )
// 	json.Unmarshal( hand_shake_bytes , &hand_shake_json )
// 	return
// }

func GetHandshakeData() ( hand_shake_json interface{} ) {
	err := json.Unmarshal( handshake_json , &hand_shake_json )
	if err != nil { panic( err ) }
	return
}