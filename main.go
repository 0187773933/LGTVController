package main

import (
	"os"
	"fmt"
	"path/filepath"
	types "github.com/0187773933/LGTVController/v1/types"
	utils "github.com/0187773933/LGTVController/v1/utils"
	lg_tv "github.com/0187773933/LGTVController/v1/controller"
)

func main() {
	config_file_path := "./config.yaml"
	if len( os.Args ) > 1 { config_file_path , _ = filepath.Abs( os.Args[ 1 ] ) }
	config := utils.GetConfig( config_file_path )
	fmt.Printf( "Loaded Config File From : %s\n" , config_file_path )

	tv := lg_tv.New( &config )

	tv.API( "power_on" )
	tv.API( "set_volume", types.Payload{
		"volume": 15 ,
	})
	tv.API( "set_mute" )
	fmt.Println( tv.API( "get_volume" ) )
}