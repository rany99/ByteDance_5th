package main

import (
	"ByteDance_5th/pkg/config"
	"ByteDance_5th/routers"
	"fmt"
)

func main() {

	//gin.SetMode(gin.ReleaseMode)

	r := routers.InitRouters()

	err := r.Run(fmt.Sprintf(":%d", config.Conf.SE.Port))

	if err != nil {
		return
	}
}
