package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"strings"
)

// Mysql mysql连接信息
type Mysql struct {
	Host      string
	Port      int
	Database  string
	Username  string
	Password  string
	Charset   string
	ParseTime bool `toml:"parse_time"`
	Loc       string
}

// Server 服务器连接信息
type Server struct {
	IP       string
	Port     int
	Database int
}

// Redis Redis连接信息
type Redis struct {
	IP       string
	Port     int
	Database int
}

// Config 配置信息
type Config struct {
	DB Mysql  `toml:"mysql"`
	RD Redis  `toml:"redis"`
	SE Server `toml:"server"`
}

var Conf Config

const tomlAddr string = "D:\\go_project\\ByteDance_5th\\pkg\\config\\config.toml"

// 初始化
func init() {
	if _, err := toml.DecodeFile(tomlAddr, &Conf); err != nil {
		panic(err)
	}
	//去除左右空格
	strings.Trim(Conf.DB.Host, " ")
	strings.Trim(Conf.RD.IP, " ")
	strings.Trim(Conf.SE.IP, " ")
	log.Println("DB.Host:", Conf.DB.Host)
	log.Println("RD.IP:", Conf.RD.IP)
	log.Println("SE.IP:", Conf.SE.IP)
}

// GetConnectionString 获取数据库连接
func GetConnectionString() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=%s",
		Conf.DB.Username,
		Conf.DB.Password,
		Conf.DB.Host,
		Conf.DB.Port,
		Conf.DB.Database,
		Conf.DB.Charset,
		Conf.DB.ParseTime,
		Conf.DB.Loc)
	log.Println(dsn)
	return dsn
}
