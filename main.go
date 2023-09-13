package main

import (
	types "github.com/0187773933/LGTVController/v1/types"
	utils "github.com/0187773933/LGTVController/v1/utils"
	lg_tv "github.com/0187773933/LGTVController/v1/controller"
)

func main() {
	config := utils.GetConfig( "config.yml" )
	tv := lg_tv.New( &config )

	tv.API( "power_on" )
	tv.API( "set_volume", types.Payload{
		"volume": 15 ,
	})
	tv.API( "set_mute" )
}