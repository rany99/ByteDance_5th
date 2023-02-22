package main

import (
	"ByteDance_5th/pkg/config"
	"ByteDance_5th/routers"
	"fmt"
	"github.com/gin-gonic/gin"
)

// 采用 Gorm 中的 AutoMigrate 自动迁移建表，无需手动建表
// 请确保 redis 服务已经开启
// 请确保 ffmpeg.exe 截图工具已经被置于 gopath 以下
// 请修改 pkg/config/config.toml 内的配置信息
// 请修改 pkg/config/config.go 中的 tomlAddr 为 config.toml 在本机的绝对路径

// go run main.go

func main() {

	// 关闭 debug_mode
	gin.SetMode(gin.ReleaseMode)

	r := routers.InitRouters()

	err := r.Run(fmt.Sprintf(":%d", config.Conf.SE.Port))

	if err != nil {
		return
	}
}
