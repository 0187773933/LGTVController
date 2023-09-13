package main

import (
	"fmt"
	utils "github.com/0187773933/LGTVController/v1/utils"
	lg_tv "github.com/0187773933/LGTVController/v1/controller"
)

func main() {

	config := utils.ParseConfig( "config.yml" )
	tv := lg_tv.New( &config )
	fmt.Println( tv )

}