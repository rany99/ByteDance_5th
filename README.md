# ByteDance_5th字节跳动第五届青训营后端进阶版大项目
ByteDance_5th ： 基于 Gin 框架与 Gorm 框架以及 Redis 实现的极简版抖音\
根据最新一版的项目说明文档实现了包括 互动方向 与 社交方向 的所有接口与功能，解决了消息记录的刷屏问题，并为用户信息新增了头像、背景、个性签名、总获赞数、作品数量等字段\
项目地址：https://github.com/303228744/ByteDance_5th\
1 项目运行
1.1 运行环境
1. Go 版本： go 1.20
2. 数据库：MySQL 8.0.26
3. Redis：3.2.100
4. ffmepg：请确保将 ffmpeg.exe 置于 GoPath/bin 下
5. 抖声：app-release.apk
1.2 配置信息
1. 请进入 ByteDance_5th\pkg\config\config.toml 修改相应的 MySQL、Redis 及 Server 端口信息
2. 请进入 ByteDance_5th\pkg\config\config.go 将 tomlAddr 修改为 config.toml 在本机的绝对路径\
1.3 项目启动
cd .\ByteDance_5th
go run main.go
