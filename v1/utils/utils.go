package utils

import (
	"os"
	"io/ioutil"
	"encoding/json"
	yaml "gopkg.in/yaml.v2"
	types "github.com/0187773933/LGTVController/v1/types"
)

func ParseConfig( file_path string ) ( result types.ConfigFile ) {
	file , _ := ioutil.ReadFile( file_path )
	error := yaml.Unmarshal( file , &result )
	if error != nil { panic( error ) }
	return
}

func GetHandshakeData() ( hand_shake_json interface{} ) {
	hand_shake_file , _ := os.Open( "./v1/utils/handshake.json" )
	defer hand_shake_file.Close()
	hand_shake_bytes , _ := ioutil.ReadAll( hand_shake_file )
	json.Unmarshal( hand_shake_bytes , &hand_shake_json )
	return
}