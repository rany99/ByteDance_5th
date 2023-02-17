package main

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/routers"
	"fmt"
)

func main() {

	//gin.SetMode(gin.ReleaseMode)

	r := routers.InitRouters()

	err := r.Run(fmt.Sprintf(":%d", common.Conf.SE.Port))

	if err != nil {
		return
	}
}
