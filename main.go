package main

import (
	"fmt"
	types "github.com/0187773933/LGTVController/v1/types"
	utils "github.com/0187773933/LGTVController/v1/utils"
	lg_tv "github.com/0187773933/LGTVController/v1/controller"
)

func main() {
	config := utils.GetConfig( "./config.yaml" )
	tv := lg_tv.New( &config )
	tv.API( "power_on" )
	tv.API( "set_volume", types.Payload{
		"volume": 15 ,
	})
	fmt.Println( tv.API( "set_mute" ) )

}