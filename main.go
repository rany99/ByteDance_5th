package main

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/routers"
	"fmt"
)

func main() {
	r := routers.InitRouters()
	err := r.Run(fmt.Sprintf(":%d", common.Conf.SE.Port)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		return
	}
}
