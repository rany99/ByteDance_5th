package main

import (
	"ByteDance_5th/config"
	"ByteDance_5th/routers"
	"fmt"
)

func main() {
	r := routers.DoushengRoutersinit()
	err := r.Run(fmt.Sprintf(":%d", config.Conf.SE.Port)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		return
	}
}
