package models

import (
	"ByteDance_5th/pkg/common"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() error {
	log.Println("init start")
	var err error
	arg := common.GetConnectionString()
	DB, err = gorm.Open(mysql.Open(arg), &gorm.Config{
		PrepareStmt:            true, //缓存预编译命令
		SkipDefaultTransaction: true, //禁用默认事务操作
	})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&UserInfo{}, &Video{}, &Comment{}, &User{})
	if err != nil {
		log.Println("init failed", err.Error())
		return err
	} else {
		log.Println("init successfully")
	}
	return err
}
